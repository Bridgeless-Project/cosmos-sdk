package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// TopUp increases the token amount of the NFT. The unlock step and the already unlocked amount stay the same,
// so the vesting time grows proportionally to the top-up.
func (k msgServer) TopUp(goctx context.Context, msg *types.MsgTopUp) (*types.MsgTopUpResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, ok := k.GetNFT(ctx, msg.Nft)
	if !ok {
		return nil, types.ErrNFTNotFound
	}

	if nft.Owner != msg.Creator {
		return nil, types.ErrNFTInvalidOwner
	}

	if msg.Amount.Denom != nft.Denom {
		return nil, errors.Wrapf(types.ErrInvalidAmount, "invalid denom %s, expected %s", msg.Amount.Denom, nft.Denom)
	}

	if nft.IsFullyUnlocked() {
		return nil, types.ErrNFTFullyUnlocked
	}

	nftAddress, _ := sdk.AccAddressFromBech32(msg.Nft)
	ownerAddress, _ := sdk.AccAddressFromBech32(msg.Creator)

	if err := k.bankKeeper.SendCoins(ctx, ownerAddress, nftAddress, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errors.Wrap(err, "failed to transfer tokens to NFT")
	}

	nft.TokenAmount = nft.TokenAmount.Add(msg.Amount)
	k.SetNFT(ctx, nft)

	return new(types.MsgTopUpResponse), nil
}
