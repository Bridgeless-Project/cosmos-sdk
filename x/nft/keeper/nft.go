package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetNFT set a specific deposit in the store from its index
func (k Keeper) SetNFT(ctx sdk.Context, v types.NFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKeyPrefix))
	b := k.cdc.MustMarshal(&v)
	store.Set(types.NFTKey(v.Address), b)
}

// GetNFT returns a NFT from its index
func (k Keeper) GetNFT(
	ctx sdk.Context,
	address string,
) (val types.NFT, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKeyPrefix))
	b := store.Get(types.NFTKey(
		address,
	))

	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNFT removes a NFT from the store
func (k Keeper) RemoveNFT(
	ctx sdk.Context,
	address string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKeyPrefix))
	store.Delete(types.NFTKey(
		address,
	))
}

// GetAllNFT returns all NFT
func (k Keeper) GetNFTs(ctx sdk.Context) (list []types.NFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NFT
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// GetNFTsWithPagination returns all NFTs with pagination
func (k Keeper) GetNFTsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]types.NFT, *query.PageResponse, error) {
	var nfts []types.NFT
	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyPrefix(types.NFTKeyPrefix))

	pageRes, err := query.Paginate(nftStore, pagination, func(key []byte, value []byte) error {
		var nft types.NFT
		k.cdc.MustUnmarshal(value, &nft)

		nfts = append(nfts, nft)
		return nil
	})

	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return nfts, pageRes, nil
}

// SetOwnerNFT set nft owner
func (k Keeper) SetOwnerNFT(ctx sdk.Context, owner, nftAddress string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	ownerBranchStore := prefix.NewStore(store, types.KeyPrefix(owner))
	data := k.cdc.MustMarshal(&types.Owner{
		Address:    owner,
		NftAddress: nftAddress,
	})
	ownerBranchStore.Set(types.NFTOwnerKey(owner), data)
}

// RemoveOwner removes a NFT from the store
func (k Keeper) RemoveOwnerNft(
	ctx sdk.Context,
	owner string,
	nftAddress string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	ownerBranchStore := prefix.NewStore(store, types.KeyPrefix(owner))
	ownerBranchStore.Delete(types.NFTKey(
		nftAddress,
	))
}

// GetAllNFTsByOwnerWithPagination returns all nfts by holder address with pagination
func (k Keeper) GetAllNFTsByOwnerWithPagination(ctx sdk.Context, ownerAddress string, pagination *query.PageRequest) ([]types.NFT, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	ownerStore := prefix.NewStore(store, types.KeyPrefix(ownerAddress))
	var nfts []types.NFT

	pageRes, err := query.Paginate(ownerStore, pagination, func(key []byte, value []byte) error {
		var owner types.Owner

		k.cdc.MustUnmarshal(value, &owner)

		nft, found := k.GetNFT(ctx, owner.Address)
		if !found {
			return types.ErrNFTNotFound
		}

		nfts = append(nfts, nft)
		return nil
	})

	if err != nil {
		return nil, pageRes, status.Error(codes.Internal, err.Error())
	}

	return nfts, pageRes, nil
}

// GetAllOwnersWithPagination returns all nft holders address with pagination
func (k Keeper) GetAllOwnersWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	var owners []string

	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
		var owner types.Owner

		k.cdc.MustUnmarshal(value, &owner)
		owners = append(owners, owner.Address)
		return nil
	})

	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return owners, pageRes, nil

}
