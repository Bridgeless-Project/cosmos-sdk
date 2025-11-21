package types

import (
	"fmt"

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
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte(ParamSuperAdminKey), &p.SuperAdmin, validateSuperAdmin),
	}
}

// Validate validates the set of params
func (p *Params) Validate() error {
	if err := validateSuperAdmin(p.SuperAdmin); err != nil {
		return errorsmod.Wrap(err, "invalid super admin address")
	}
	return nil
}

func validateSuperAdmin(i interface{}) error {
	admin, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid x/%s module admin address: %s", ModuleName, err.Error())
	}

	return nil
}
