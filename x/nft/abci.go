package nft

import (
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
		if ctx.BlockHeight() < int64(nft.StartVestingBlock) {
			continue
		}

		if ctx.BlockTime().Unix()-nft.LastVestingTime < nft.VestingPeriod {
			continue
		}

		if nft.VestingCounter >= nft.VestingPeriodsCount {
			continue
		}

		nft.AvailableToWithdraw = nft.AvailableToWithdraw.Add(sdk.NewCoin(nft.Denom, nft.RewardPerPeriod.Amount))
		nft.VestingCounter++
		nft.LastVestingTime = ctx.BlockTime().Unix()

		k.SetNFT(ctx, nft)
	}

}
