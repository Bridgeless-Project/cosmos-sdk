package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Undelegate(goctx context.Context, request *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Address)
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

	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)

	_, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	_, err = k.stakingKeeper.Undelegate(ctx, nftAddress, valAddr, sdk.NewDecCoinFromCoin(request.Amount).Amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to undelegate tokens")
	}

	return new(types.MsgUndelegateResponse), nil
}

func (k msgServer) Delegate(goctx context.Context, request *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)
	nft, ok := k.GetNFT(ctx, request.Address)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	valAddr, _ := sdk.ValAddressFromBech32(request.Validator)
	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrTokenIsDelegated
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	_, err := k.stakingKeeper.Delegate(ctx, nftAddress, request.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		k.Logger(ctx).Error("failed to delegate", "error", err)
		return nil, types.ErrFailedToDelegate
	}

	return new(types.MsgDelegateResponse), nil

}

func (k msgServer) Send(goctx context.Context, request *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Address)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrTokenIsDelegated
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	k.RemoveOwnerNft(ctx, nft.Owner, nft.Address)
	nft.Owner = request.Recipient

	k.SetOwnerNFT(ctx, nft.Owner, nft.Address)
	k.SetNFT(ctx, nft)
	return new(types.MsgSendResponse), nil

}

func (k msgServer) Withdraw(goctx context.Context, request *types.MsgWithdrawal) (*types.MsgWithdrawalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Address)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)
	ownerAddress, _ := sdk.AccAddressFromBech32(request.Creator)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrTokenIsDelegated
	}

	balance := k.bankKeeper.GetAllBalances(ctx, nftAddress).AmountOf(nft.Denom)
	amount := sdk.MinInt(nft.AvailableToWithdraw.Amount, balance)

	err := k.DistributeToAddress(ctx, sdk.NewCoins(sdk.NewCoin(nft.Denom, amount)), ownerAddress, nftAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to distribute to address")
	}

	nft.AvailableToWithdraw = nft.AvailableToWithdraw.Sub(sdk.NewCoin(nft.Denom, amount))
	k.SetNFT(ctx, nft)

	return new(types.MsgWithdrawalResponse), nil
}

func (k msgServer) Redelegate(goctx context.Context, request *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, request.Address)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != request.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)
	validatorSrcAddress, _ := sdk.ValAddressFromBech32(request.ValidatorSrc)
	validatorNEwAddress, _ := sdk.ValAddressFromBech32(request.ValidatorNew)

	_, err := k.stakingKeeper.BeginRedelegation(ctx, nftAddress, validatorSrcAddress, validatorNEwAddress, sdk.NewDecCoinFromCoin(request.Amount).Amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin redelegation")
	}

	return new(types.MsgRedelegateResponse), nil
}
