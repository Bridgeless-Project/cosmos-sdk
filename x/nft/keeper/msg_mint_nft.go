package keeper

import (
	"context"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	accumulatorKeeper "github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	accumulatortypes "github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (m msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := m.Keeper.GetParams(sdkCtx)

	if !m.IsModuleAdmin(sdkCtx, msg.Creator) {
		return nil, sdkerrors.Wrapf(errors.ErrUnauthorized, "invalid NFT creator %s", msg.Creator)
	}

	vestingPeriod := big.NewInt(params.VestingPeriod)
	totalVestingTime := big.NewInt(params.TotalVestingTime)

	vestingPeriodsCount := big.NewInt(0).Div(totalVestingTime, vestingPeriod)
	nftTokenAnount, ok := big.NewInt(0).SetString(params.NftTokenAmount, 10)
	if !ok {
		return nil, sdkerrors.Wrap(errors.ErrInvalidRequest, "invalid NFT cost")
	}

	vestingRewardPerPeriod := new(big.Int).Div(nftTokenAnount, vestingPeriodsCount)

	nft, sequence, err := m.CreateNft(sdkCtx, msg.Owner, msg.StartVestingBlock, vestingRewardPerPeriod, msg.NftMetadataUri)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create NFT")
	}

	nftAddress, err := sdk.AccAddressFromBech32(nft.Address)
	if err != nil {
		return nil, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid NFT address %s", nft.Address)
	}

	nftPoolAddress := accumulatorKeeper.GetPoolAddress(accumulatortypes.NFTPoolName)
	poolBalances := m.bankKeeper.GetAllBalances(sdkCtx, nftPoolAddress)
	ok, balance := poolBalances.Find(m.GetBondDenom(sdkCtx))
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidBalance, "balance not found")
	}

	if balance.Amount.LT(sdk.NewInt(nftTokenAnount.Int64())) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient pool balance, NFT cost is %d, balance: %d",
			params.NftTokenAmount, balance.Amount.Int64())
	}

	coinsToDistribute := sdk.NewCoin(params.BondDenom, sdk.NewIntFromBigInt(nftTokenAnount))

	if err = m.accumulatorKeeper.DistributeToAccount(
		sdkCtx,
		accumulatortypes.NFTPoolName,
		sdk.NewCoins(coinsToDistribute),
		nftAddress,
	); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to mint coins to NFT")
	}

	m.SetNFT(sdkCtx, *nft)
	params.NftSequence = sequence
	m.SetParams(sdkCtx, params)

	return &types.MsgMintResponse{}, nil
}
