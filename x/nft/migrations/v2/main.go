package v2

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info(fmt.Sprintf("Performing v12.1.8-rc2 %s module migrations", types.ModuleName))

	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.NFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var nft types.NFT
		cdc.MustUnmarshal(iterator.Value(), &nft)

		nftOwnerStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
		ownerBranchStore := prefix.NewStore(nftOwnerStore, types.KeyPrefix(nft.Owner))
		ownerBranchStore.Delete(types.NFTOwnerKey(
			nft.Owner,
		))

		data := cdc.MustMarshal(&types.Owner{
			Address:    nft.Owner,
			NftAddress: nft.Address,
		})
		ownerBranchStore.Set(types.NFTOwnerKey(nft.Address), data)
	}

	return nil
}
