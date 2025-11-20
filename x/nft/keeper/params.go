package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
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

func (k Keeper) GetVestingTime(ctx sdk.Context) (vestingTime int64) {
	k.paramstore.Get(ctx, []byte(types.ParamsNftTotalVestingTimeKey), &vestingTime)
	return
}

func (k Keeper) GetVestingPeriod(ctx sdk.Context) (vestingPeriod int64) {
	k.paramstore.Get(ctx, []byte(types.ParamsNftVestingPeriodKey), &vestingPeriod)

	return
}

func (k Keeper) GetVestingPeriodsCount(ctx sdk.Context) (vestingPeriodCount uint64) {
	k.paramstore.Get(ctx, []byte(types.ParamVestingCountKey), &vestingPeriodCount)
	return
}

func (k Keeper) GetNFTToeknAmount(ctx sdk.Context) (nftTokenAmount string) {
	k.paramstore.Get(ctx, []byte(types.ParamNftTokenAmountKey), &nftTokenAmount)
	return
}
