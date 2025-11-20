package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) AddAdmin(goctx context.Context, req *types.MsgAddAdmin) (*types.MsgAddAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if m.GetParams(ctx).SuperAdmin != req.Creator {
		return nil, types.ErrForbidden
	}

	if _, ok := m.GetAdmin(ctx, req.Address); ok {
		return nil, types.ErrAdminExists
	}
	newAdmin := types.Admin{
		Address:             req.Address,
		VestingPeriod:       req.VestingPeriod,
		RewardPerPeriod:     req.RewardPerPeriod,
		VestingPeriodsCount: req.VestingPeriodsCount,
		VestingCounter:      0,
		LastVestingTime:     0,
		Denom:               req.Denom,
	}

	m.SetAdmin(ctx, newAdmin)

	return new(types.MsgAddAdminResponse), nil
}

func (m msgServer) BurnTokens(goctx context.Context, req *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if m.GetParams(ctx).SuperAdmin != req.Sender {
		return nil, types.ErrForbidden
	}

	err := m.BurnTokensFromPool(ctx, req.PoolName, req.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "burning tokens from pool")
	}

	return &types.MsgBurnTokensResponse{}, nil
}
