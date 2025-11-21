package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"golang.org/x/net/context"
)

type (
	Keeper interface {
		Logger(c sdk.Context) log.Logger
		GetParams(ctx sdk.Context) (params types.Params)
		SetParams(ctx sdk.Context, params types.Params)
		DistributeToModule(ctx sdk.Context, pool string, amount sdk.Coins, receiverModule string) error
		DistributeToAccount(ctx sdk.Context, pool string, amount sdk.Coins, receiver sdk.AccAddress) error
		SetAdmin(ctx sdk.Context, admin types.Admin)
		GetAdmin(ctx sdk.Context, address string) (val types.Admin, found bool)
		RemoveAdmin(ctx sdk.Context, address string)
		GetAllAdmins(ctx sdk.Context) []types.Admin
		GetAdmins(ctx context.Context, admins *types.QueryAdmins) (*types.QueryAdminsResponse, error)
		GetAdminByAddress(ctx context.Context, admin *types.QueryAdminByAddress) (*types.QueryAdminByAddressResponse, error)
		BurnTokensFromPool(ctx sdk.Context, pool string, amount sdk.Coin) error
		GetSuperAdmin(ctx sdk.Context) (superAdmin string)
		SetSuperAdmin(ctx sdk.Context, superAdmin string)
	}

	BaseKeeper struct {
		cdc             codec.BinaryCodec
		storeKey        storetypes.StoreKey
		memKey          storetypes.StoreKey
		paramstore      paramtypes.Subspace
		bankKeeper      bankkeeper.Keeper
		ak              accountKeeper.AccountKeeper
		lastVestingTime time.Time
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak accountKeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return BaseKeeper{
		bankKeeper: bankKeeper,
		cdc:        cdc,
		paramstore: ps,
		storeKey:   storeKey,
		ak:         ak,
		memKey:     memKey,
	}
}

func (k BaseKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func GetPoolAddress(poolName string) sdk.AccAddress {
	return address.Derive(authtypes.NewModuleAddress(types.ModuleName), []byte(poolName))
}
