package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

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
