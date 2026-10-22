package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (k msgServer) Withdraw(goctx context.Context, msg *types.MsgWithdrawal) (*types.MsgWithdrawalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, msg.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != msg.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	nftAddress, _ := sdk.AccAddressFromBech32(msg.Nft)
	ownerAddress, _ := sdk.AccAddressFromBech32(msg.Creator)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrNFTIsDelegated
	}

	if msg.Amount.Denom != nft.Denom {
		return nil, errors.Wrapf(types.ErrInvalidAmount, "invalid denom %s, expected %s", msg.Amount.Denom, nft.Denom)
	}

	amount := msg.Amount.Amount
	balance := k.bankKeeper.GetBalance(ctx, nftAddress, nft.Denom).Amount

	// fully unlocked: any part of the balance, including rewards and donations, can be withdrawn.
	// otherwise the limit depends on the token amount, not on the live balance, so rewards and donations do not increase it.
	limit := balance
	if !nft.IsFullyUnlocked() {
		limit = sdk.MinInt(nft.WithdrawableBeforeUnlock(), balance)
	}

	if amount.GT(limit) {
		return nil, errors.Wrapf(types.ErrInvalidAmount, "amount %s exceeds the withdrawable limit %s", amount, limit)
	}

	if !nft.IsFullyUnlocked() {
		nft.WithdrawnAmount = nft.WithdrawnAmount.Add(msg.Amount)
	}

	err := k.DistributeToAddress(ctx, sdk.NewCoins(sdk.NewCoin(nft.Denom, amount)), ownerAddress, nftAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to distribute to address")
	}

	k.SetNFT(ctx, nft)

	return new(types.MsgWithdrawalResponse), nil
}
