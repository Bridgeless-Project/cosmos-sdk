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

	currentBlockHeight := ctx.BlockHeader().Height

	for _, nft := range nfts {
		if nft.StartVestingBlock == 0 {
			nft.StartVestingBlock = currentBlockHeight
		}

		if currentBlockHeight < nft.StartVestingBlock {
			continue
		}

		if nft.VestingPeriodsCount >= int64(params.VestingPeriodsLimit) {
			continue
		}

		// if not full period passed since last vesting skip the nft
		if currentBlockHeight-nft.LastVestingBlock < params.VestingPeriod {
			continue
		}

		// if vesting time has passed skip the nft
		if currentBlockHeight-nft.StartVestingBlock >= params.TotalVestingTime {
			continue
		}

		currentVestingPeriod := currentBlockHeight - nft.LastVestingBlock
		if nft.LastVestingBlock == 0 {
			currentVestingPeriod = params.VestingPeriod
		}

		passedPeriods := big.NewInt(0).Div(big.NewInt(currentVestingPeriod), big.NewInt(params.VestingPeriod))

		if nft.VestingPeriodsCount+passedPeriods.Int64() > int64(params.VestingPeriodsLimit) {
			passedPeriods = big.NewInt(int64(params.VestingPeriodsLimit) - nft.VestingPeriodsCount)
		}

		nft.AvailableToWithdraw = nft.AvailableToWithdraw.Add(sdk.NewCoin(
			nft.Denom,
			nft.RewardPerPeriod.Amount.Mul(sdk.NewInt(passedPeriods.Int64())),
		),
		)

		nft.VestingPeriodsCount += passedPeriods.Int64()
		nft.LastVestingBlock = currentBlockHeight

		k.SetNFT(ctx, nft)
	}

	params.BatchIndex++
	k.SetParams(ctx, params)
}
