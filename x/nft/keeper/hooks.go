package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
