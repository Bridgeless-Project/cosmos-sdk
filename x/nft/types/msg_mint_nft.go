package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMint = "mint"
)

var _ sdk.Msg = &MsgMintRequest{}

func NewMsgMint(creator, address string, startVestingBlock int64) *MsgMintRequest {
	return &MsgMintRequest{
		Creator:           creator,
		Owner:             address,
		StartVestingBlock: startVestingBlock,
	}
}

func (msg *MsgMintRequest) Route() string {
	return RouterKey
}

func (msg *MsgMintRequest) Type() string {
	return TypeMsgMint
}

func (msg *MsgMintRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if msg.StartVestingBlock < 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "negative start vesting block")
	}

	return nil
}
