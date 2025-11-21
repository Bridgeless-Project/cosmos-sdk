package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

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
