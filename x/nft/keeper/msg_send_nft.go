package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (k msgServer) Send(goctx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := k.validateIsNFT(ctx, msg.Recipient); err != nil {
		return nil, errors.Wrap(err, "recipient address is NFT")
	}

	nft, ok := k.GetNFT(ctx, msg.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	nftAddress, _ := sdk.AccAddressFromBech32(msg.Nft)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrNFTIsDelegated
	}

	if nft.Owner != msg.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	k.RemoveOwnerNft(ctx, nft.Owner, nft.Address)
	nft.Owner = msg.Recipient

	k.SetOwnerNFT(ctx, nft.Owner, nft.Address)
	k.SetNFT(ctx, nft)
	return new(types.MsgSendResponse), nil

}
