package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgStakingRewardsWithdrawal = "staking-rewards-withdrawal"
)

var _ sdk.Msg = &MsgStakingWithdrawalRewards{}

func NewMsgStakingWithdrawalRewards(creator, validator, nft string) *MsgStakingWithdrawalRewards {
	return &MsgStakingWithdrawalRewards{
		Creator:   creator,
		Validator: validator,
		Nft:       nft,
	}
}

func (msg *MsgStakingWithdrawalRewards) Route() string {
	return RouterKey
}

func (msg *MsgStakingWithdrawalRewards) Type() string {
	return TypeMsgStakingRewardsWithdrawal
}

func (msg *MsgStakingWithdrawalRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStakingWithdrawalRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStakingWithdrawalRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Nft)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid NFT address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}
