package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (k msgServer) Redelegate(goctx context.Context, request *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Nft)
	validatorSrcAddress, _ := sdk.ValAddressFromBech32(request.ValidatorSrc)
	validatorNEwAddress, _ := sdk.ValAddressFromBech32(request.ValidatorNew)
	share, err := sdk.NewDecFromStr(request.Share)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode share")
	}

	_, err = k.stakingKeeper.BeginRedelegation(
		ctx,
		nftAddress,
		validatorSrcAddress,
		validatorNEwAddress,
		share,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin redelegation")
	}

	return new(types.MsgRedelegateResponse), nil
}

func (k msgServer) validateIsNFT(ctx sdk.Context, address string) error {
	_, ok := k.GetNFT(ctx, address)
	if !ok {
		return nil
	}

	return types.ErrAddressISNFT
}

func (k msgServer) Undelegate(goctx context.Context, request *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	valAddr, err := sdk.ValAddressFromBech32(request.Validator)
	if err != nil {
		return nil, err
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Nft)

	_, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	share, err := sdk.NewDecFromStr(request.Share)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode share")
	}
	_, err = k.stakingKeeper.Undelegate(ctx, nftAddress, valAddr, share)
	if err != nil {
		return nil, errors.Wrap(err, "failed to undelegate tokens")
	}

	return new(types.MsgUndelegateResponse), nil
}

func (k msgServer) Delegate(goctx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)
	valAddr, _ := sdk.ValAddressFromBech32(msg.Validator)
	nftAddress, _ := sdk.AccAddressFromBech32(msg.Nft)
	delegator, _ := sdk.AccAddressFromBech32(msg.Creator)
	if err := k.DelegateNFT(ctx,
		nftAddress,
		delegator,
		valAddr,
		msg.Amount.Amount,
		false,
	); err != nil {
		return nil, errors.Wrap(err, "failed to delegate NFT")
	}

	return new(types.MsgDelegateResponse), nil
}
