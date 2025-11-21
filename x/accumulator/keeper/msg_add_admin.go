package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

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
