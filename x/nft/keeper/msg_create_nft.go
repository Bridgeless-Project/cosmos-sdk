package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (m msgServer) Create(ctx context.Context, request *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if !m.IsModuleAdmin(sdkCtx, request.Creator) {
		return nil, sdkerrors.Wrapf(errors.ErrUnauthorized, "invalid NFT creator %s", request.Creator)
	}

	ownerAddress, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid owner address %s", request.Owner)
	}

	ownerBalances := m.bankKeeper.GetAllBalances(sdkCtx, ownerAddress)
	ok, balance := ownerBalances.Find(m.BondDenom(sdkCtx))
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidBalance, "balance not found")
	}

	if balance.Amount.LT(sdk.NewInt(NFTCost.Int64())) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient balance, NFT cost is %d, balance: %d",
			NFTCost.Int64(), balance.Amount.Int64())
	}

	if err = m.bankKeeper.SendCoinsFromAccountToModule(
		sdkCtx,
		ownerAddress, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(m.BondDenom(sdkCtx), sdk.NewInt(NFTCost.Int64()))),
	); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to send coins to module account")
	}

	nft, err := m.CreateNft(sdkCtx, request.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create NFT")
	}

	m.SetNFT(sdkCtx, *nft)

	return &types.MsgCreateResponse{}, nil
}
