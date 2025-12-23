package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oldTypes "github.com/cosmos/cosmos-sdk/x/nft/migrations/v3/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, params types.Params) error {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.NFTKeyPrefix))

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	nftTokensValue, ok := sdk.NewIntFromString(params.NftTokenAmount)
	if !ok {
		return fmt.Errorf("invalid nftTokenAmount")
	}

	for ; iterator.Valid(); iterator.Next() {
		var oldNft oldTypes.NFT
		cdc.MustUnmarshal(iterator.Value(), &oldNft)

		var newNft = types.NFT{
			Address:             oldNft.Address,
			Owner:               oldNft.Owner,
			Uri:                 oldNft.Uri,
			RewardPerPeriod:     oldNft.RewardPerPeriod,
			VestingPeriodsCount: oldNft.VestingPeriodsCount,
			AvailableToWithdraw: oldNft.AvailableToWithdraw,
			LastVestingBlock:    oldNft.LastVestingBlock,
			Denom:               oldNft.Denom,
			StartVestingBlock:   oldNft.StartVestingBlock,
			VestingPeriodsLimit: params.VestingPeriodsLimit,
			IsFrozen:            false,
			TotalVestingTime:    params.TotalVestingTime,
			TokenAmount:         sdk.NewCoin(oldNft.Denom, nftTokensValue),
		}

		nftOwnerStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
		ownerBranchStore := prefix.NewStore(nftOwnerStore, types.KeyPrefix(oldNft.Owner))
		ownerBranchStore.Delete(types.NFTOwnerKey(
			newNft.Owner,
		))

		store.Set(types.NFTKey(oldNft.Address), cdc.MustMarshal(&newNft))
	}

	return nil
}
