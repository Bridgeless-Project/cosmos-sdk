package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v043 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v043"
	v046 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v046"
	v04626 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v046_26"
	v04628 "github.com/cosmos/cosmos-sdk/x/staking/migrations/v046_28"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v043.MigrateStore(ctx, m.keeper.storeKey)
}

// Migrate2to3 migrates x/staking state from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v046.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.paramstore)
}

// Migrate2to3 migrates x/staking state from consensus version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v04626.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate4to5 migrates x/staking state from consensus version 4 to 5.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v04628.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.keeper.paramstore)
}
