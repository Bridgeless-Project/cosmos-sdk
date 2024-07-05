package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (k BaseKeeper) DistributeValidatorsPool(ctx sdk.Context, amount sdk.Coins) error {
	return k.distributeTokens(ctx, types.ValidatorPoolName, authtypes.NewModuleAddress(minttypes.ModuleName), amount)
}

func (k BaseKeeper) distributeTokens(ctx sdk.Context, fromPool string, receiver sdk.AccAddress, amount sdk.Coins) error {
	poolAddress := GetPoolAddress(fromPool)
	if poolAddress == nil {
		return types.ErrInvalidPool
	}

	return k.sendFromAddressToAddress(ctx, poolAddress, receiver, amount)
}

func (k BaseKeeper) sendFromAddressToAddress(ctx sdk.Context, poolAddress, receiverAddress sdk.AccAddress, amount sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		poolAddress,
		types.ModuleName,
		amount,
	)

	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		receiverAddress,
		amount,
	)

	if err != nil {
		return errors.Wrap(err, "sending native coins to address")
	}

	return nil
}
