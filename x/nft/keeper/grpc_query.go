package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	ctx := sdk.UnwrapSDKContext(c)
	nfts, page, err := k.GetNFTsWithPagination(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllNFTsResponse{
		Nft:        nfts,
		Pagination: page,
	}, nil
}

func (k Keeper) GetAllNFTsByOwner(c context.Context, req *types.QueryAllNFTsByOwner) (*types.QueryAllNFTsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	nfts, page, err := k.GetAllNFTsByOwnerWithPagination(sdk.UnwrapSDKContext(c), req.Owner, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNFTsByOwnerResponse{
		Nft:        nfts,
		Pagination: page,
	}, nil
}

func (k Keeper) GetAllOwners(c context.Context, req *types.QueryAllOwners) (*types.QueryAllOwnersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	owners, page, err := k.GetAllOwnersWithPagination(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOwnersResponse{
		Owner:      owners,
		Pagination: page,
	}, nil
}
