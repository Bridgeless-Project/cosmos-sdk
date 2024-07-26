package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GetNFTByAddress(goctx context.Context, req *types.QueryNFTByAddress) (*types.QueryNFTByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, found := k.GetNFT(ctx, req.Address)
	if !found {
		return nil, types.ErrNFTNotFound
	}

	return &types.QueryNFTByAddressResponse{
		Nft: &nft,
	}, nil
}

func (k Keeper) GetAllNFTs(c context.Context, req *types.QueryAllNFTs) (*types.QueryAllNFTsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nfts []types.NFT
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyPrefix(types.NFTKeyPrefix))

	pageRes, err := query.Paginate(nftStore, req.Pagination, func(key []byte, value []byte) error {
		var nft types.NFT
		k.cdc.MustUnmarshal(value, &nft)

		nfts = append(nfts, nft)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNFTsResponse{Nft: nfts, Pagination: pageRes}, nil
}

func (k Keeper) GetAllNFTsByOwner(c context.Context, req *types.QueryAllNFTsByOwner) (*types.QueryAllNFTsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var nfts []types.NFT
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	ownerStore := prefix.NewStore(store, types.KeyPrefix(req.Owner))

	pageRes, err := query.Paginate(ownerStore, req.Pagination, func(key []byte, value []byte) error {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNFTsByOwnerResponse{Nft: nfts, Pagination: pageRes}, nil
}

func (k Keeper) GetAllOwners(c context.Context, req *types.QueryAllOwners) (*types.QueryAllOwnersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var owners []string
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NFTByOwnerKeyPrefix))
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var owner types.Owner

		k.cdc.MustUnmarshal(value, &owner)
		owners = append(owners, owner.Address)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOwnersResponse{Owner: owners, Pagination: pageRes}, nil
}
