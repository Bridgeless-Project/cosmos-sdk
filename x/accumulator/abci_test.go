package accumulator_test

import (
	"github.com/cosmos/cosmos-sdk/x/accumulator"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestVerifyVesting(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{
		Time: time.Now(),
	})

	admins := app.AccumulatorKeeper.GetAllAdmins(ctx)
	require.Empty(t, admins)

	adminPool, _ := sdk.AccAddressFromBech32("cosmos1mr9ydjyr7qjuuu39qur9756pe5l8dpqlfhxd7efd2wxj58zdm9dqw07l5t")

	app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, adminPool, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000))))
	newAdmin := types.Admin{
		Address:             "cosmos1k6vrturm2249y705ygsr9lz79y30aq90sltvyr",
		VestingPeriod:       1,
		RewardPerPeriod:     sdk.NewCoin("stake", sdk.NewInt(1000)),
		VestingPeriodsCount: 10,
		VestingCounter:      0,
		LastVestingTime:     0,
		Denom:               "stake",
	}

	adminAddress, _ := sdk.AccAddressFromBech32(newAdmin.Address)
	app.AccumulatorKeeper.SetAdmin(ctx, newAdmin)

	for i := 0; i < int(newAdmin.VestingPeriodsCount); i++ {
		accumulator.EndBlocker(ctx, app.AccumulatorKeeper)

		balance := app.BankKeeper.GetBalance(ctx, adminAddress, "stake")
		require.Equal(t, balance.Amount, newAdmin.RewardPerPeriod.Amount.MulRaw(int64(i+1)))

		time.Sleep(time.Duration(newAdmin.VestingPeriod) * time.Second)
		ctx = app.BaseApp.NewContext(true, tmproto.Header{
			Time: time.Now(),
		})
	}

	balance := app.BankKeeper.GetBalance(ctx, adminAddress, "stake")
	require.Equal(t, balance.Amount, newAdmin.RewardPerPeriod.Amount.MulRaw(newAdmin.VestingPeriodsCount))
}

func TestVerifyVestingForListOfAdmins(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{
		Time: time.Now(),
	})

	admins := app.AccumulatorKeeper.GetAllAdmins(ctx)
	require.Empty(t, admins)

	adminPool, _ := sdk.AccAddressFromBech32("cosmos1mr9ydjyr7qjuuu39qur9756pe5l8dpqlfhxd7efd2wxj58zdm9dqw07l5t")
	app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, adminPool, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000))))

	admin1 := types.Admin{
		Address:             "cosmos1k6vrturm2249y705ygsr9lz79y30aq90sltvyr",
		VestingPeriod:       1,
		RewardPerPeriod:     sdk.NewCoin("stake", sdk.NewInt(1000)),
		VestingPeriodsCount: 10,
		VestingCounter:      0,
		LastVestingTime:     0,
		Denom:               "stake",
	}

	admin2 := types.Admin{
		Address:             "cosmos19tfnk0s4m095828tlm22c5q0awqf7nd97szd0z",
		VestingPeriod:       3,
		RewardPerPeriod:     sdk.NewCoin("stake", sdk.NewInt(123)),
		VestingPeriodsCount: 12,
		VestingCounter:      0,
		LastVestingTime:     0,
		Denom:               "stake",
	}
	admin3 := types.Admin{
		Address:             "cosmos1k9ulywn8fhupem5cdcyp3rlrx2mg2n9larm2ha",
		VestingPeriod:       6,
		RewardPerPeriod:     sdk.NewCoin("stake", sdk.NewInt(323)),
		VestingPeriodsCount: 8,
		VestingCounter:      0,
		LastVestingTime:     0,
		Denom:               "stake",
	}

	app.AccumulatorKeeper.SetAdmin(ctx, admin1)
	app.AccumulatorKeeper.SetAdmin(ctx, admin2)
	app.AccumulatorKeeper.SetAdmin(ctx, admin3)

	for i := 0; i < 50; i++ {
		accumulator.EndBlocker(ctx, app.AccumulatorKeeper)
		time.Sleep(time.Duration(1) * time.Second)
		ctx = app.BaseApp.NewContext(true, tmproto.Header{
			Time: time.Now(),
		})
	}

	for _, admin := range app.AccumulatorKeeper.GetAllAdmins(ctx) {
		address, _ := sdk.AccAddressFromBech32(admin.Address)
		balance := app.BankKeeper.GetBalance(ctx, address, "stake")
		require.Equal(t, balance.Amount, admin.RewardPerPeriod.Amount.MulRaw(admin.VestingPeriodsCount))

	}
}
