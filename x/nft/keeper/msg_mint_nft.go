package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	accumulatortypes "github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (m msgServer) Mint(ctx context.Context, request *types.MsgMintRequest) (*types.MsgMintResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := m.Keeper.GetParams(sdkCtx)

	if !m.IsModuleAdmin(sdkCtx, request.Creator) {
		return nil, sdkerrors.Wrapf(errors.ErrUnauthorized, "invalid NFT creator %s", request.Creator)
	}

	nft, err := m.CreateNft(sdkCtx, request.Owner, request.StartVestingBlock)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create NFT")
	}

	nftAddress, err := sdk.AccAddressFromBech32(nft.Address)
	if err != nil {
		return nil, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid NFT address %s", nft.Address)
	}

	nftPoolAddress := m.accumulatorKeeper.GetPoolAddress(accumulatortypes.NFTPoolName)
	poolBalances := m.bankKeeper.GetAllBalances(sdkCtx, nftPoolAddress)
	ok, balance := poolBalances.Find(m.GetBondDenom(sdkCtx))
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidBalance, "balance not found")
	}

	if balance.Amount.LT(sdk.NewInt(NFTCost.Int64())) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient pool balance, NFT cost is %d, balance: %d",
			NFTCost.Int64(), balance.Amount.Int64())
	}

	if err = m.accumulatorKeeper.DistributeToAccount(
		sdkCtx,
		accumulatortypes.NFTPoolName,
		sdk.NewCoins(
			sdk.NewCoin(m.GetBondDenom(sdkCtx), sdk.NewInt(NFTCost.Int64())),
		),
		nftAddress,
	); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to mint coins to NFT")
	}

	m.SetNFT(sdkCtx, *nft)
	params.NftSequence++
	m.SetParams(sdkCtx, params)

	return &types.MsgMintResponse{}, nil
}
