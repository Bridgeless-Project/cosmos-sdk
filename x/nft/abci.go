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
	nfts, _, _ := k.GetNFTsWithPagination(
		ctx,
		&query.PageRequest{
			Limit:  params.BatchSize,
			Offset: params.BatchIndex * params.BatchSize,
		},
	)

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

		// if nft is fully unlocked skip it
		if nft.IsFullyUnlocked() {
			continue
		}

		if currentBlockHeight < nft.StartVestingBlock {
			continue
		}

		// if not full period passed since last vesting skip the nft
		if currentBlockHeight-nft.LastVestingBlock < params.VestingPeriod {
			continue
		}

		currentVestingPeriod := currentBlockHeight - nft.LastVestingBlock
		if nft.LastVestingBlock == 0 {
			currentVestingPeriod = params.VestingPeriod
		}

		passedPeriods := sdk.NewInt(currentVestingPeriod / params.VestingPeriod)

		// unlock a fixed step per period, the last period unlocks only the remainder so the unlock ends exactly at 100%
		remainder := nft.TokenAmount.Amount.Sub(nft.UnlockedAmount.Amount)
		toUnlock := sdk.MinInt(nft.RewardPerPeriod.Amount.Mul(passedPeriods), remainder)

		nft.UnlockedAmount = nft.UnlockedAmount.Add(sdk.NewCoin(nft.Denom, toUnlock))
		nft.LastVestingBlock = currentBlockHeight

		k.SetNFT(ctx, nft)
	}

	params.BatchIndex++
	k.SetParams(ctx, params)
}
