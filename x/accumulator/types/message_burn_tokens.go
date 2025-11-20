package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeBurnTokens = "burn_tokens"
)

var _ sdk.Msg = &MsgBurnTokens{}

func NewMsgBurnTokens(sender, poolName string, amount sdk.Coin) *MsgBurnTokens {
	return &MsgBurnTokens{
		Sender:   sender,
		PoolName: poolName,
		Amount:   amount,
	}
}

func (msg *MsgBurnTokens) Route() string {
	return RouterKey
}

func (msg *MsgBurnTokens) Type() string {
	return TypeBurnTokens
}

func (msg *MsgBurnTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	err = msg.Amount.Validate()
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid amount (%s)", err)
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	if msg.PoolName == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "pool name cannot be empty")
	}

	return nil
}
