package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple bank hooks, all hook functions are run in array sequence
var _ BankHooks = &MultiBankHooks{}

type MultiBankHooks []BankHooks

func NewMultiBankHooks(hooks ...BankHooks) MultiBankHooks {
	return hooks
}

func (h MultiBankHooks) BeforeSendTokenToAddress(ctx sdk.Context, sender, receiver sdk.Address, coins sdk.Coins) error {
	for i := range h {
		if err := h[i].BeforeSendTokenToAddress(ctx, sender, receiver, coins); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiBankHooks) AfterSendTokenToAddress(ctx sdk.Context, receiver sdk.Address, coins sdk.Coins) error {
	for i := range h {
		if err := h[i].AfterSendTokenToAddress(ctx, receiver, coins); err != nil {
			return err
		}
	}

	return nil
}
