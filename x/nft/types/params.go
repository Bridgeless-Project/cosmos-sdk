package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(moduleAdmin, bondDenom, prefix string) Params {
	return Params{
		ModuleAdmin: moduleAdmin,
		BondDenom:   bondDenom,
		Prefix:      prefix,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams("", sdk.DefaultBondDenom, "prefix")
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte(ParamModuleAdminKey), &p.ModuleAdmin, validateModuleAdmin),
		paramtypes.NewParamSetPair([]byte(ParamBondDenomKey), &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair([]byte(ParamPrefixKey), &p.Prefix, validatePrefix),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateModuleAdmin(p.ModuleAdmin); err != nil {
		return errorsmod.Wrap(err, "invalid module admin")
	}

	if err := validateBondDenom(p.BondDenom); err != nil {
		return errorsmod.Wrap(err, "invalid bond denom")
	}

	if err := validatePrefix(p.Prefix); err != nil {
		return errorsmod.Wrap(err, "invalid prefix")
	}

	return nil
}

func validateModuleAdmin(i interface{}) error {
	adm, ok := i.(string)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(adm)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid x/%s module admin address: %s", ModuleName, err.Error())
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	bondDenom, ok := i.(string)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter type: %T", i)
	}

	if strings.TrimSpace(bondDenom) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "empty denomination")
	}

	return nil
}

func validatePrefix(i interface{}) error {
	prefix, ok := i.(string)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter type: %T", i)
	}

	if strings.TrimSpace(prefix) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "empty prefix")
	}

	return nil
}
