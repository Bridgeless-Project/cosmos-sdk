package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v3types "github.com/cosmos/cosmos-sdk/x/nft/migrations/v3/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// MigrateStore rewrites every NFT from the legacy layout (available_to_withdraw, vesting_periods_count, ...)
// to the new one (unlocked_amount, withdrawn_amount).
//
// The legacy available_to_withdraw was decremented on every withdrawal, so the cumulative unlocked amount is
// restored from vesting_periods_count * reward_per_period (capped by token_amount) and the already withdrawn
// part is the difference. This keeps unlocked - withdrawn equal to the legacy available_to_withdraw.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info(fmt.Sprintf("Performing %s module v3 migrations", types.ModuleName))

	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.NFTKeyPrefix))

	// collect first: the store must not be written while iterating
	type entry struct {
		key []byte
		nft types.NFT
	}
	var entries []entry

	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	for ; iterator.Valid(); iterator.Next() {
		var old v3types.NFT
		if err := cdc.Unmarshal(iterator.Value(), &old); err != nil {
			iterator.Close()
			return fmt.Errorf("failed to unmarshal legacy NFT %s: %w", string(iterator.Key()), err)
		}

		entries = append(entries, entry{
			key: append([]byte{}, iterator.Key()...),
			nft: convert(old),
		})
	}
	iterator.Close()

	for _, e := range entries {
		store.Set(e.key, cdc.MustMarshal(&e.nft))
	}

	return nil
}

func convert(old v3types.NFT) types.NFT {
	denom := old.Denom
	zero := sdk.ZeroInt()

	tokenAmount := old.TokenAmount.Amount
	if tokenAmount.IsNil() {
		tokenAmount = zero
	}
	available := old.AvailableToWithdraw.Amount
	if available.IsNil() {
		available = zero
	}
	reward := old.RewardPerPeriod.Amount
	if reward.IsNil() {
		reward = zero
	}

	unlocked := sdk.MinInt(reward.MulRaw(old.VestingPeriodsCount), tokenAmount)
	// available can not exceed what was unlocked; guard against inconsistent legacy state
	if available.GT(unlocked) {
		unlocked = sdk.MinInt(available, tokenAmount)
		available = unlocked
	}
	withdrawn := unlocked.Sub(available)

	return types.NFT{
		Address:           old.Address,
		Owner:             old.Owner,
		Uri:               old.Uri,
		RewardPerPeriod:   sdk.NewCoin(denom, reward),
		LastVestingBlock:  old.LastVestingBlock,
		Denom:             denom,
		StartVestingBlock: old.StartVestingBlock,
		TokenAmount:       sdk.NewCoin(denom, tokenAmount),
		UnlockedAmount:    sdk.NewCoin(denom, unlocked),
		WithdrawnAmount:   sdk.NewCoin(denom, withdrawn),
	}
}
