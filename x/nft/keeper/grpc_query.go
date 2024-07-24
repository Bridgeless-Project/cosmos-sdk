package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GetNFTByAddress(goctx context.Context, reqeust *types.QueryNFTByAddress) (*types.QueryNFTByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	nft, found := k.GetNFT(ctx, reqeust.Address)
	if !found {
		return nil, types.ErrNFTNotFound
	}

	return &types.QueryNFTByAddressResponse{
		Nft: &nft,
	}, nil
}
