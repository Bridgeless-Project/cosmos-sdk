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

func (h MultiBankHooks) BeforeSendTokenToAddress(ctx sdk.Context, receiver sdk.Address) error {
	for i := range h {
		if err := h[i].BeforeSendTokenToAddress(ctx, receiver); err != nil {
			return err
		}
	}

	return nil
}
