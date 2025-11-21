package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

// GetParams get all parameters as types.Params
func (k BaseKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k BaseKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k BaseKeeper) GetSuperAdmin(ctx sdk.Context) (superAdmin string) {
	k.paramstore.Get(ctx, []byte(types.ParamSuperAdminKey), &superAdmin)
	return
}

func (k BaseKeeper) SetSuperAdmin(ctx sdk.Context, superAdmin string) {
	k.paramstore.Set(ctx, []byte(types.ParamSuperAdminKey), []byte(superAdmin))
}
