package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/tendermint/tendermint/libs/log"
	"golang.org/x/net/context"
	"time"
)

type (
	Keeper interface {
		Logger(c sdk.Context) log.Logger
		GetParams(c sdk.Context) types.Params
		SetParams(c sdk.Context, params types.Params)
		Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error)
		DistributeTokens(ctx sdk.Context, fromPool string, isSentToModule bool, amount sdk.Coins, receiverModule string, receiverAddress *sdk.AccAddress) error
	}

	BaseKeeper struct {
		cdc             codec.BinaryCodec
		storeKey        storetypes.StoreKey
		memKey          storetypes.StoreKey
		bankKeeper      bankkeeper.Keeper
		ak              accountKeeper.AccountKeeper
		lastVestingTime time.Time
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ak accountKeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) *BaseKeeper {
	return &BaseKeeper{
		bankKeeper: bankKeeper,
		cdc:        cdc,
		storeKey:   storeKey,
		ak:         ak,
		memKey:     memKey,
	}
}

func (k BaseKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func GetPoolAddress(poolName string) sdk.AccAddress {
	return sdk.AccAddress(address.Derive(authtypes.NewModuleAddress(types.ModuleName), []byte(poolName)))
}
