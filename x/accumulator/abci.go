package accumulator

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"time"
)

// update vesting state for each admin
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	for _, admin := range k.GetAllAdmins(ctx) {
		if ctx.BlockTime().Unix()-admin.LastVestingTime < admin.VestingPeriod {
			return
		}

		if admin.VestingCounter >= admin.VestingPeriodsCount {
			return
		}

		address, _ := sdk.AccAddressFromBech32(admin.Address)
		err := k.DistributeToAccount(ctx, types.AdminPoolName, sdk.NewCoins(sdk.NewCoin(admin.Denom, admin.RewardPerPeriod.Amount)), address)
		if err != nil {
			return
		}

		admin.VestingCounter++
		admin.LastVestingTime = ctx.BlockTime().Unix()

		k.SetAdmin(ctx, admin)
	}

}
