package keeper

import (
	"context"
	"math/big"
	"time"

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

	nftCost := m.GetNFTCost(sdkCtx)
	vestingTime := m.GetVestingTime(sdkCtx)
	vestingPeriod := m.GetVestingPeriod(sdkCtx)

	period := big.NewInt(vestingPeriod)
	vestTime := big.NewInt(vestingTime)

	vestingPeriodsCount := big.NewInt(0).Div(vestTime, period)
	cost, ok := big.NewInt(0).SetString(nftCost, 10)
	if !ok {
		return nil, sdkerrors.Wrap(errors.ErrInvalidRequest, "invalid NFT cost")
	}

	vestingReward := new(big.Int).Div(cost, vestingPeriodsCount)

	nft, err := m.CreateNft(sdkCtx, request.Owner, time.Now().Unix(), vestingReward)
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

	if balance.Amount.LT(sdk.NewInt(cost.Int64())) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient pool balance, NFT cost is %d, balance: %d",
			NFTCost.Int64(), balance.Amount.Int64())
	}

	coinsToDistribute := sdk.NewCoin(m.GetBondDenom(sdkCtx), sdk.NewIntFromBigInt(cost))

	if err = m.accumulatorKeeper.DistributeToAccount(
		sdkCtx,
		accumulatortypes.NFTPoolName,
		sdk.NewCoins(coinsToDistribute),
		nftAddress,
	); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to mint coins to NFT")
	}

	m.SetNFT(sdkCtx, *nft)
	params.NftSequence++
	m.SetParams(sdkCtx, params)

	return &types.MsgMintResponse{}, nil
}
