package types

import (
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtype "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	TypeMsgBecameValidator = "became_validator"
)

var _ sdk.Msg = &MsgBecomeValidator{}

// NewMsgBecomeValidator creates a new MsgBecomeValidator instance.
// Delegator address and validator address are the same.
func NewMsgBecomeValidator(
	valAddr sdk.ValAddress,
	pubKey cryptotypes.PubKey, //nolint:interfacer
	description stakingtype.Description,
	commission stakingtype.CommissionRates,
	minSelfDelegation math.Int,
	nftAddresses []string,
) (*MsgBecomeValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgBecomeValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr).String(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
		NftAddresses:      nftAddresses,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgBecomeValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgBecomeValidator) Type() string { return TypeMsgBecameValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgBecomeValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	addrs := []sdk.AccAddress{delegator}
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)

	valAccAddr := sdk.AccAddress(valAddr)
	if !delegator.Equals(valAccAddr) {
		addrs = append(addrs, valAccAddr)
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgBecomeValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBecomeValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures both non-empty and valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.Pubkey == nil {
		return stakingtype.ErrEmptyValidatorPubKey
	}

	if msg.Description == (stakingtype.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtype.CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err = msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	for _, nftAddr := range msg.NftAddresses {
		if _, err = sdk.AccAddressFromBech32(nftAddr); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, nftAddr)
		}
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgBecomeValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}
