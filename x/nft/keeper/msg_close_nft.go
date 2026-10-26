package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// Close sends the whole balance of the NFT to its owner without waiting for the end of vesting and removes the NFT.
func (k msgServer) Close(goctx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if !k.IsModuleAdmin(ctx, msg.Creator) {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "invalid NFT closer %s", msg.Creator)
	}

	nft, ok := k.GetNFT(ctx, msg.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	nftAddress, err := sdk.AccAddressFromBech32(nft.Address)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid NFT address %s", nft.Address)
	}

	ownerAddress, err := sdk.AccAddressFromBech32(nft.Owner)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid NFT owner %s", nft.Owner)
	}

	if k.HasStake(ctx, nftAddress) {
		return nil, errors.Wrap(types.ErrNFTIsDelegated, "can not close NFT: tokens are delegated or being unbonded")
	}

	balances := k.bankKeeper.GetAllBalances(ctx, nftAddress)
	if !balances.IsZero() {
		if err := k.DistributeToAddress(ctx, balances, ownerAddress, nftAddress); err != nil {
			return nil, errors.Wrap(err, "failed to distribute to address")
		}
	}

	k.RemoveOwnerNft(ctx, nft.Owner, nft.Address)
	k.RemoveNFT(ctx, nft.Address)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNFTClosed,
			sdk.NewAttribute(types.AttributeKeyNFT, nft.Address),
			sdk.NewAttribute(types.AttributeKeyOwner, nft.Owner),
			sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
		),
	)

	return new(types.MsgCloseResponse), nil
}
