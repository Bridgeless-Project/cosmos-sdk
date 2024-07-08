package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (k BaseKeeper) DistributeValidatorsPool(ctx sdk.Context, amount sdk.Coins) error {
	return k.distributeTokens(ctx, types.ValidatorPoolName, minttypes.ModuleName, amount)
}

func (k BaseKeeper) distributeTokens(ctx sdk.Context, fromPool string, receiverModule string, amount sdk.Coins) error {
	poolAddress := GetPoolAddress(fromPool)
	if poolAddress == nil {
		return types.ErrInvalidPool
	}

	return k.sendFromAddressToModule(ctx, poolAddress, receiverModule, amount)
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
