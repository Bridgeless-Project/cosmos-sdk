package keeper

import (
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ banktypes.BankHooks = Hooks{}

// Hooks creates new nft hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// BeforeSendTokenToAddress handles cases when recipient is a NFT
func (h Hooks) BeforeSendTokenToAddress(_ sdk.Context, _, _ sdk.Address, _ sdk.Coins) error {
	// this hook is ignored
	return nil
}

// AfterSendTokenToAddress updates vesting params for NFT
func (h Hooks) AfterSendTokenToAddress(ctx sdk.Context, receiver sdk.Address, amt sdk.Coins) error {
	params := h.k.GetParams(ctx)

	nft, found := h.k.GetNFT(ctx, receiver.String())
	if !found {
		return nil
	}

	// get count of additional periods
	additionalPeriods := math.Ceil(float64(amt.AmountOf(nft.Denom).Int64()) / float64(nft.RewardPerPeriod.Amount.Int64()))

	nft.VestingPeriodsLimit += int64(additionalPeriods)
	ok, tokenToAdd := amt.Find(nft.Denom)
	if !ok {
		return sdkerrors.Wrapf(types.ErrInvalidAmount, "token to add with denom %s doesn't exist", nft.Denom)
	}

	nft.TokenAmount = nft.TokenAmount.Add(tokenToAdd)

	blockDistanceFromLastVesting := ctx.BlockHeight() - nft.LastVestingBlock

	// timeToRestartVesting includes time to get to batch where nft will be processed and covers the time was spent after last vesting block
	timeToRestartVesting := h.k.GetNftBatchBlockDistance(ctx, nft.Address) + blockDistanceFromLastVesting

	// adding the time of vesting
	nft.TotalVestingTime += int64(additionalPeriods)*params.VestingPeriod + timeToRestartVesting

	h.k.SetNFT(ctx, nft)

	return nil
}
