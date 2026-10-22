package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgTopUp = "top_up"
)

var _ sdk.Msg = &MsgTopUp{}

func NewMsgTopUp(creator, nft string, amount sdk.Coin) *MsgTopUp {
	return &MsgTopUp{
		Creator: creator,
		Nft:     nft,
		Amount:  amount,
	}
}

func (msg *MsgTopUp) Route() string {
	return RouterKey
}

func (msg *MsgTopUp) Type() string {
	return TypeMsgTopUp
}

func (msg *MsgTopUp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTopUp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTopUp) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Nft); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid nft address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid amount %s", msg.Amount)
	}

	return nil
}
