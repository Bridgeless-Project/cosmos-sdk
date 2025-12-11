package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMint = "mint"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(creator, address, metadataUri string, startVestingBlock int64) *MsgMint {
	return &MsgMint{
		Creator:           creator,
		Owner:             address,
		StartVestingBlock: startVestingBlock,
		NftMetadataUri:    metadataUri,
	}
}

func (msg *MsgMint) Route() string {
	return RouterKey
}

func (msg *MsgMint) Type() string {
	return TypeMsgMint
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if strings.TrimSpace(msg.NftMetadataUri) == "" {
		return sdkerrors.Wrap(ErrEmptyNftMetadataUri, "nft metadata uri is required")
	}

	if msg.StartVestingBlock < 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "negative start vesting block")
	}

	return nil
}
