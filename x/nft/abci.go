package nft

import (
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// update vesting state for each nft
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	params := k.GetParams(ctx)
	nfts, _, _ := k.GetNFTsWithPagination(ctx, &query.PageRequest{Limit: params.BatchSize, Offset: params.BatchIndex * params.BatchSize})
	if len(nfts) == 0 {
		params.BatchIndex = 0
		k.SetParams(ctx, params)
		return
	}

	for _, nft := range nfts {

		if ctx.BlockTime().Unix() < nft.StartVestingTime {
			continue
		}

		if nft.VestingCounter >= int64(params.VestingPeriodsCount) {
			continue
		}

		// if not full period passed since last vesting skip the nft
		if ctx.BlockTime().Unix()-nft.LastVestingTime < params.VestingPeriod {
			continue
		}

		// if vesting time has passed skip the nft
		if ctx.BlockTime().Unix()-nft.StartVestingTime >= params.VestingTime {
			continue
		}

		currentVestingTime := ctx.BlockTime().Unix() - nft.LastVestingTime
		if nft.LastVestingTime == 0 {
			currentVestingTime = params.VestingPeriod
		}

		passedPeriods := big.NewInt(0).Div(big.NewInt(currentVestingTime), big.NewInt(params.VestingPeriod))

		if nft.VestingCounter+passedPeriods.Int64() > int64(params.VestingPeriodsCount) {
			passedPeriods = big.NewInt(int64(params.VestingPeriodsCount) - nft.VestingPeriodsCount)
		}

		nft.AvailableToWithdraw = nft.AvailableToWithdraw.Add(sdk.NewCoin(
			nft.Denom,
			nft.RewardPerPeriod.Amount.Mul(sdk.NewInt(passedPeriods.Int64())),
		),
		)

		nft.VestingCounter += passedPeriods.Int64()
		nft.LastVestingTime = ctx.BlockTime().Unix()

		k.SetNFT(ctx, nft)
	}

	params.BatchIndex++
	k.SetParams(ctx, params)
}
