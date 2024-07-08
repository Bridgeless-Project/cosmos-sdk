package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

func (k BaseKeeper) DistributeTokens(ctx sdk.Context, fromPool string, isSentToModule bool, amount sdk.Coins, receiverModule string, receiverAddress *sdk.AccAddress) error {
	poolAddress := GetPoolAddress(fromPool)
	if poolAddress == nil {
		return types.ErrInvalidPool
	}

	if isSentToModule {
		return k.sendFromAddressToModule(ctx, poolAddress, receiverModule, amount)
	}

	if receiverAddress == nil {
		return types.ErrInvalidReceiver
	}
	return k.sendFromAddressToAddress(ctx, poolAddress, *receiverAddress, amount)

}

func (k BaseKeeper) sendFromAddressToModule(ctx sdk.Context, poolAddress sdk.AccAddress, receiverAddress string, amount sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		poolAddress,
		receiverAddress,
		amount,
	)

	if err != nil {
		err = errors.Wrap(err, "sending native coins to address")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (k BaseKeeper) sendFromAddressToAddress(ctx sdk.Context, poolAddress sdk.AccAddress, receiverAddress sdk.AccAddress, amount sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		poolAddress,
		types.ModuleName,
		amount,
	)

	if err != nil {
		err = errors.Wrap(err, "sending native coins to module")
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
		err = errors.Wrap(err, "sending native coins to module")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	return nil
}
