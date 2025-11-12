package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

var NFTCost = new(big.Int).Mul(big.NewInt(1280000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

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

func (k Keeper) GetModuleAdmin(ctx sdk.Context) (adminAddress string) {
	k.paramstore.Get(ctx, []byte(types.ParamModuleAdminKey), &adminAddress)
	return
}

func (k Keeper) GetBondDenom(ctx sdk.Context) (bondDenom string) {
	k.paramstore.Get(ctx, []byte(types.ParamBondDenomKey), &bondDenom)
	return
}

func (k Keeper) GetPrefix(ctx sdk.Context) (prefix string) {
	k.paramstore.Get(ctx, []byte(types.ParamPrefixKey), &prefix)
	return
}

func (k Keeper) GetNftSequence(ctx sdk.Context) (seq uint64) {
	k.paramstore.Get(ctx, []byte(types.ParamNftSequenceKey), &seq)
	return
}

func (k Keeper) GetVestingPeriod(ctx sdk.Context) (vestingPeriod uint64) {
	k.paramstore.Get(ctx, []byte(types.ParamsNftVestingPeriodKey), &vestingPeriod)
	return
}

func (k Keeper) GetVestingPeriodsCount(ctx sdk.Context) (vestingPeriodCount uint64) {
	k.paramstore.Get(ctx, []byte(types.ParamVestingCountKey), &vestingPeriodCount)
	return
}

func (k Keeper) GetNFTCost(ctx sdk.Context) (nftCost uint64) {
	k.paramstore.Get(ctx, []byte(types.ParamNftCostKey), &nftCost)
	return
}
