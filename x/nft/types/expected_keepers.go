package types

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins // Methods imported from bank should be defined here
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AccumulatorKeeper interface {
	DistributeToAccount(ctx sdk.Context, pool string, amount sdk.Coins, receiver sdk.AccAddress) error
}

type StakingKeeper interface {
	GetAllDelegatorDelegations(ctx sdk.Context, address sdk.AccAddress) []stakingtypes.Delegation
	Delegate(
		ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares sdk.Dec, err error)
	Undelegate(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
	) (time.Time, error)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	BeginRedelegation(
		ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
	) (completionTime time.Time, err error)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	BondDenom(ctx sdk.Context) (res string)
	MinDelegationAmount(ctx sdk.Context) (res string)
	SetValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	SetValidatorByConsAddr(ctx sdk.Context, validator stakingtypes.Validator) error
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
	SetNewValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error
}
