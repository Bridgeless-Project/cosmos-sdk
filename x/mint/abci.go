package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// fetch stored minter & params
	params := k.GetParams(ctx)

	// if block reqard is zero, skip any logic
	if params.BlockReward.Amount.IsZero() {
		k.Logger(ctx).Info("The block reward is zero")
		return
	}

	// validate halving params
	if uint64(ctx.BlockHeight())%params.HalvingBlocks == 0 && params.CurrentHalvingPeriod < params.MaxHalvingPeriods {
		// halving the rewards
		params.BlockReward.Amount = params.BlockReward.Amount.Quo(sdk.NewInt(2))
		params.CurrentHalvingPeriod++

		// set zero reward if max halving periods reached
		if params.CurrentHalvingPeriod == params.MaxHalvingPeriods {
			params.BlockReward.Amount = sdk.ZeroInt()
		}

		// set the updated params
		k.SetParams(ctx, params)
	}

	// mint coins, update supply
	mintedCoins := sdk.NewCoins(params.BlockReward)
	err := k.SendFromAccumulator(ctx, mintedCoins)
	if err != nil {
		k.Logger(ctx).Error("failed to send tokens from accumulator")
		return
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		k.Logger(ctx).Error("failed to collect fees")
		return
	}

	if params.BlockReward.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(params.BlockReward.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, params.BlockReward.Amount.String()),
		),
	)
}
