package accumulator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// This line is used by starport scaffolding # genesis/module/init
	for _, admin := range genState.Admins {
		k.SetAdmin(ctx, admin)
	}

	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Admins = k.GetAllAdmins(ctx)

	// this line is used by starport scaffolding # genesis/module/export
	return genesis
}
