package keeper

import (
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
func (h Hooks) BeforeSendTokenToAddress(ctx sdk.Context, sender, receiver sdk.Address, amt sdk.Coins) error {
	nft, found := h.k.GetNFT(ctx, receiver.String())
	if !found {
		return nil
	}

	if !sender.Equals(sdk.MustAccAddressFromBech32(nft.Owner)) {
		return sdkerrors.Wrap(types.ErrNFTInvalidOwner, "sender is not the owner of the NFT")
	}

	// validate that user can send only multiple of reward per period
	if !amt.AmountOf(nft.Denom).Mod(nft.RewardPerPeriod.Amount).IsZero() {
		return sdkerrors.Wrap(types.ErrInvalidAmount, "amount is not a multiple of reward per period")
	}

	return nil
}

// AfterSendTokenToAddress updates vesting params for NFT
func (h Hooks) AfterSendTokenToAddress(ctx sdk.Context, receiver sdk.Address, amt sdk.Coins) error {
	nft, found := h.k.GetNFT(ctx, receiver.String())
	if !found {
		return nil
	}

	// get count of additional periods
	additionalPeriods := amt.AmountOf(nft.Denom).Quo(nft.RewardPerPeriod.Amount)
	nft.VestingPeriodsCount += additionalPeriods.Int64()
	h.k.SetNFT(ctx, nft)

	return nil
}
