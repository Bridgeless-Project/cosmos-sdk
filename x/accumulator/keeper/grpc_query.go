package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = BaseKeeper{}

func (k BaseKeeper) GetAdmins(c context.Context, req *types.QueryAdmins) (*types.QueryAdminsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	admins, page, err := k.GetAllAdminsWithPagination(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAdminsResponse{
		Admins:     admins,
		Pagination: page,
	}, nil
}

func (k BaseKeeper) GetAdminByAddress(c context.Context, req *types.QueryAdminByAddress) (*types.QueryAdminByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	admin, ok := k.GetAdmin(ctx, req.Address)
	if !ok {
		return nil, status.Error(codes.Internal, types.ErrAdminNotFound.Error())
	}

	return &types.QueryAdminByAddressResponse{
		Admin: admin,
	}, nil
}
