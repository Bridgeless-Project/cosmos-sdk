package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ banktypes.BankHooks = Hooks{}

// Hooks creates new nft hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// BeforeSendTokenToAddress handles that receiver address isn't a NFT
func (h Hooks) BeforeSendTokenToAddress(ctx sdk.Context, receiver sdk.Address) error {
	_, found := h.k.GetNFT(ctx, receiver.String())
	if found {
		return sdkerrors.Wrap(types.ErrAddressISNFT, "receiver is a NFT address")
	}

	// TODO some additional business logic here
	return nil
}
