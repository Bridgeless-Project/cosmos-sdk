package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeAddAdmin = "send"
)

var _ sdk.Msg = &MsgAddAdmin{}

func NewMsgAddAdmin(creator, denom, address string, rewardPerPeriod sdk.Coin, vestingPeriodsCount int64, vestingPeriod int64) *MsgAddAdmin {
	return &MsgAddAdmin{
		Creator:             creator,
		RewardPerPeriod:     rewardPerPeriod,
		Address:             address,
		VestingPeriodsCount: vestingPeriodsCount,
		VestingPeriod:       vestingPeriod,
		Denom:               denom,
	}
}

func (msg *MsgAddAdmin) Route() string {
	return RouterKey
}

func (msg *MsgAddAdmin) Type() string {
	return TypeAddAdmin
}

func (msg *MsgAddAdmin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddAdmin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}

	if msg.VestingPeriod <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid vesting period")
	}

	if msg.VestingPeriodsCount <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid vesting period count")
	}

	return nil
}
