package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	accumulatorKeeper "github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	accumulatortypes "github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	paramSpace        paramtypes.Subspace
	stakingKeeper     types.StakingKeeper
	bankKeeper        types.BankKeeper
	accumulatorKeeper accumulatorKeeper.Keeper
	feeCollectorName  string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper, acc accumulatorKeeper.Keeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:               cdc,
		storeKey:          key,
		paramSpace:        paramSpace,
		stakingKeeper:     sk,
		bankKeeper:        bk,
		feeCollectorName:  feeCollectorName,
		accumulatorKeeper: acc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) math.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.stakingKeeper.BondedRatio(ctx)
}

func (k Keeper) SendFromAccumulator(ctx context.Context, moduleName string, amount sdk.Coins) error {
	fmt.Println("SendFromModuleToAddressViaAccumulator")
	request := accumulatortypes.DistributeTokensRequest{
		Amount:         amount,
		ModuleNameFrom: accumulatortypes.ModuleName,
		ModuleNameTo:   moduleName,
	}

	msgServer := accumulatorKeeper.NewMsgServerImpl(k.accumulatorKeeper)
	_, err := msgServer.DistributeTokens(ctx, &request)
	if err != nil {
		err = errors.Wrap(err, "failed to call accumulator module")
		k.Logger(sdk.UnwrapSDKContext(ctx)).Error(err.Error())
		return err
	}

	return nil
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
