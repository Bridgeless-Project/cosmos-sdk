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

func (m msgServer) Mint(ctx context.Context, request *types.MsgMintRequest) (*types.MsgMintResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := m.Keeper.GetParams(sdkCtx)

	if !m.IsModuleAdmin(sdkCtx, request.Creator) {
		return nil, sdkerrors.Wrapf(errors.ErrUnauthorized, "invalid NFT creator %s", request.Creator)
	}

	vestingPeriod := big.NewInt(params.VestingPeriod)
	totalVestingTime := big.NewInt(params.TotalVestingTime)

	vestingPeriodsCount := big.NewInt(0).Div(totalVestingTime, vestingPeriod)
	nftTokenAnount, ok := big.NewInt(0).SetString(params.NftTokenAmount, 10)
	if !ok {
		return nil, sdkerrors.Wrap(errors.ErrInvalidRequest, "invalid NFT cost")
	}

	vestingRewardPerPeriod := new(big.Int).Div(nftTokenAnount, vestingPeriodsCount)

	nft, sequence, err := m.CreateNft(sdkCtx, request.Owner, request.StartVestingBlock, vestingRewardPerPeriod)
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

	if balance.Amount.LT(sdk.NewIntFromBigInt(nftTokenAnount)) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient pool balance, NFT cost is %s, balance: %s",
			params.NftTokenAmount, balance.Amount.String())
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
