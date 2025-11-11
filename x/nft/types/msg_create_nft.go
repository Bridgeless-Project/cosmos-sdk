package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreate = "create"
)

var _ sdk.Msg = &MsgCreateRequest{}

func NewMsgCreate(creator, address string) *MsgCreateRequest {
	return &MsgCreateRequest{
		Creator: creator,
		Owner:   address,
	}
}

func (msg *MsgCreateRequest) Route() string {
	return RouterKey
}

func (msg *MsgCreateRequest) Type() string {
	return TypeMsgCreate
}

func (msg *MsgCreateRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return nil
}
