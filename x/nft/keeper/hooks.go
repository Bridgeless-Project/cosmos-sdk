package keeper

import (
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ banktypes.BankHooks = Hooks{}

// Hooks creates new nft hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// BeforeSendTokenToAddress handles cases when recipient is a NFT
func (h Hooks) BeforeSendTokenToAddress(_ sdk.Context, _, _ sdk.Address, _ sdk.Coins) error {
	// this hook is ignored
	return nil
}

// AfterSendTokenToAddress updates vesting params for NFT
func (h Hooks) AfterSendTokenToAddress(ctx sdk.Context, receiver sdk.Address, amt sdk.Coins) error {
	params := h.k.GetParams(ctx)

	nft, found := h.k.GetNFT(ctx, receiver.String())
	if !found {
		return nil
	}

	// get count of additional periods
	periodsToAddFloat, _ := new(big.Int).Quo(
		new(big.Int).SetInt64(amt.AmountOf(nft.Denom).Int64()),
		new(big.Int).SetInt64(nft.RewardPerPeriod.Amount.Int64()),
	).Float64()

	additionalPeriods := math.Ceil(periodsToAddFloat)

	nft.VestingPeriodsLimit += int64(additionalPeriods)

	nft.TokenAmount = nft.TokenAmount.Add(sdk.NewCoin(nft.Denom, amt.AmountOf(nft.Denom)))

	blockDistanceFromLastVesting := ctx.BlockHeight() - nft.LastVestingBlock

	// timeToRestartVesting includes time to get to batch where nft will be processed and covers the time was spent after last vesting block
	timeToRestartVesting := h.k.GetNftBatchBlockDistance(ctx, nft.Address) + uint64(blockDistanceFromLastVesting)

	// adding the time of vesting
	nft.TotalVestingTime += int64(additionalPeriods)*params.VestingPeriod + int64(timeToRestartVesting)

	h.k.SetNFT(ctx, nft)

	return nil
}
