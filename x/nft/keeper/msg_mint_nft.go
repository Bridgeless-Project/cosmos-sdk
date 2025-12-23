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

func (m msgServer) Mint(goctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)
	params := m.Keeper.GetParams(ctx)

	if !m.IsModuleAdmin(ctx, msg.Creator) {
		return nil, sdkerrors.Wrapf(errors.ErrUnauthorized, "invalid NFT creator %s", msg.Creator)
	}

	vestingPeriod := big.NewInt(params.VestingPeriod)
	totalVestingTime := big.NewInt(params.TotalVestingTime)

	vestingPeriodsCount := new(big.Int).Div(totalVestingTime, vestingPeriod)
	nftTokenAnount, ok := new(big.Int).SetString(params.NftTokenAmount, 10)
	if !ok {
		return nil, sdkerrors.Wrap(errors.ErrInvalidRequest, "invalid NFT cost")
	}

	vestingRewardPerPeriod := new(big.Int).Div(nftTokenAnount, vestingPeriodsCount)

	totalNftTokenAmount := sdk.NewCoin(params.BondDenom, sdk.NewIntFromBigInt(nftTokenAnount))

	nft, sequence, err := m.CreateNft(ctx, msg.Owner, msg.StartVestingBlock, vestingRewardPerPeriod, msg.NftMetadataUri, totalNftTokenAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create NFT")
	}

	nftAddress, err := sdk.AccAddressFromBech32(nft.Address)
	if err != nil {
		return nil, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid NFT address %s", nft.Address)
	}

	nftPoolAddress := accumulatorKeeper.GetPoolAddress(accumulatortypes.NFTPoolName)
	poolBalances := m.bankKeeper.GetAllBalances(ctx, nftPoolAddress)
	ok, balance := poolBalances.Find(m.GetBondDenom(ctx))
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidBalance, "balance not found")
	}

	if balance.Amount.LT(sdk.NewIntFromBigInt(nftTokenAnount)) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBalance, "insufficient pool balance, NFT cost is %s, balance: %s",
			params.NftTokenAmount, balance.Amount.String())
	}

	coinsToDistribute := sdk.NewCoin(params.BondDenom, sdk.NewIntFromBigInt(nftTokenAnount))

	err = m.accumulatorKeeper.DistributeToAccount(
		ctx,
		accumulatortypes.NFTPoolName,
		sdk.NewCoins(coinsToDistribute),
		nftAddress,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to mint coins to NFT")
	}

	m.Keeper.SetOwnerNFT(ctx, nft.Owner, nft.Address)
	m.Keeper.SetNFT(ctx, *nft)

	params.NftSequence = sequence
	m.Keeper.SetParams(ctx, params)

	return new(types.MsgMintResponse), nil
}
