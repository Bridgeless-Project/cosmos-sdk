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
	nfts, _, _ := k.GetNFTsWithPagination(ctx, &query.PageRequest{Limit: params.BatchSize, Offset: params.BatchIndex * 100})
	params.BatchIndex++
	k.SetParams(ctx, params)

	for _, nft := range nfts {
		if ctx.BlockTime().Unix() < nft.StartVestingTime {
			continue
		}

		if params.VestingPeriodsCount != 0 && nft.VestingCounter >= int64(params.VestingPeriodsCount) {
			continue
		}

		// if not full period passed since last vesting skip the nft
		if ctx.BlockTime().Unix()-nft.LastVestingTime < params.VestingPeriod {
			continue
		}

		// if vesting time has passed skip the nft
		if ctx.BlockTime().Unix()-nft.StartVestingTime > params.VestingTime {
			continue
		}

		if nft.VestingCounter >= nft.VestingPeriodsCount {
			continue
		}

		currentVestingTime := ctx.BlockTime().Unix() - nft.LastVestingTime
		passedPeriods := big.NewInt(0).Div(big.NewInt(currentVestingTime), big.NewInt(params.VestingPeriod))

		nft.AvailableToWithdraw = nft.AvailableToWithdraw.Add(sdk.NewCoin(
			nft.Denom,
			nft.RewardPerPeriod.Amount.Mul(sdk.NewInt(passedPeriods.Int64())),
		),
		)

		nft.VestingCounter++
		nft.LastVestingTime = ctx.BlockTime().Unix()

		k.SetNFT(ctx, nft)
	}

}
