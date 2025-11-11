package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) ModuleAdmin(ctx sdk.Context) (adminAddress string) {
	k.paramstore.Get(ctx, []byte(types.ParamModuleAdminKey), &adminAddress)
	return
}

func (k Keeper) BondDenom(ctx sdk.Context) (bondDenom string) {
	k.paramstore.Get(ctx, []byte(types.ParamBondDenomKey), &bondDenom)
	return
}

func (k Keeper) Prefix(ctx sdk.Context) (prefix string) {
	k.paramstore.Get(ctx, []byte(types.ParamPrefixKey), &prefix)
	return
}
