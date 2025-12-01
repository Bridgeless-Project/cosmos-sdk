package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (k msgServer) Send(goctx context.Context, request *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := k.validateIsNFT(ctx, request.Recipient); err != nil {
		return nil, errors.Wrap(err, "recipient address is NFT")
	}

	nft, ok := k.GetNFT(ctx, request.Address)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	nftAddress, _ := sdk.AccAddressFromBech32(request.Address)
	if k.IsDelegated(ctx, nftAddress) {
		return nil, types.ErrNFTIsAlreadyDelegated
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
