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
func NewParams(
	moduleAdmin,
	bondDenom,
	prefix string,
	sequence uint64,
	vestingPeriodsLimit int64,
	batchSize,
	batchIndex uint64,
	vestingTime,
	vestingPeriod int64,
	nftCost string) Params {
	return Params{
		ModuleAdmin:         moduleAdmin,
		BondDenom:           bondDenom,
		Prefix:              prefix,
		NftSequence:         sequence,
		TotalVestingTime:    vestingTime,
		VestingPeriod:       vestingPeriod,
		NftTokenAmount:      nftCost,
		VestingPeriodsLimit: vestingPeriodsLimit,
		BatchSize:           batchSize,
		BatchIndex:          batchIndex,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"",
		sdk.DefaultBondDenom,
		"prefix",
		0,
		0,
		1,
		0,
		1,
		1,
		"1")
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte(ParamModuleAdminKey), &p.ModuleAdmin, validateModuleAdmin),
		paramtypes.NewParamSetPair([]byte(ParamBondDenomKey), &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair([]byte(ParamPrefixKey), &p.Prefix, validatePrefix),
		paramtypes.NewParamSetPair([]byte(ParamNftSequenceKey), &p.NftSequence, validateNftSequence),
		paramtypes.NewParamSetPair([]byte(ParamsNftTotalVestingTimeKey), &p.TotalVestingTime, validateVestingTime),
		paramtypes.NewParamSetPair([]byte(ParamVestingPeriodsLimitKey), &p.VestingPeriodsLimit, validateVestingCount),
		paramtypes.NewParamSetPair([]byte(ParamNftTokenAmountKey), &p.NftTokenAmount, validateNftCost),
		paramtypes.NewParamSetPair([]byte(ParamBatchIndexKey), &p.BatchIndex, validateBatchIndex),
		paramtypes.NewParamSetPair([]byte(ParamBatchSizeKey), &p.BatchSize, validateBatchSize),
		paramtypes.NewParamSetPair([]byte(ParamsNftVestingPeriodKey), &p.VestingPeriod, validateVestingPeriod),
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

	if err := validateNftSequence(p.NftSequence); err != nil {
		return errorsmod.Wrap(err, "invalid nft sequence")
	}

	if err := validateVestingTime(p.TotalVestingTime); err != nil {
		return errorsmod.Wrap(err, "invalid vesting time")
	}

	if err := validateNftCost(p.NftTokenAmount); err != nil {
		return errorsmod.Wrap(err, "invalid nft token amount")
	}

	if err := validateVestingCount(p.VestingPeriodsLimit); err != nil {
		return errorsmod.Wrap(err, "invalid vesting period")
	}

	if err := validateBatchSize(p.BatchSize); err != nil {
		return errorsmod.Wrap(err, "invalid batch size")
	}

	if err := validateBatchIndex(p.BatchIndex); err != nil {
		return errorsmod.Wrap(err, "invalid batch index")
	}

	if err := validateVestingPeriod(p.VestingPeriod); err != nil {
		return errorsmod.Wrap(err, "invalid vesting period")
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
		return errorsmod.Wrap(ErrInvalidBondDenom, "empty denomination")
	}

	return nil
}

func validatePrefix(i interface{}) error {
	prefix, ok := i.(string)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid parameter type: %T", i)
	}

	if strings.TrimSpace(prefix) == "" {
		return errorsmod.Wrap(ErrInvalidPrefix, "empty prefix")
	}

	return nil
}

func validateNftSequence(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid nft sequence type: %T", i)
	}

	return nil
}

func validateVestingTime(i interface{}) error {
	time, ok := i.(int64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid vesting time type: %T", i)
	}

	if time == 0 {
		return errorsmod.Wrap(ErrInvalidVestingPeriod, "zero vesting time")
	}

	return nil
}

func validateVestingCount(i interface{}) error {
	val, ok := i.(int64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid vesting limit type: %T", i)
	}

	if val < 0 {
		return errorsmod.Wrap(ErrInvalidVestingPeriodsLimit, "negative vesting periods limit")
	}

	return nil
}

func validateNftCost(i interface{}) error {
	d, ok := sdk.NewIntFromString(i.(string))
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid nft cost type")
	}

	if d.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidNftCost, "nft cost amount must be positive: %s", d)
	}

	if d.IsZero() {
		return errorsmod.Wrap(ErrInvalidNftCost, "nft cost amount is zero")
	}

	return nil
}

func validateBatchSize(i interface{}) error {
	size, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid batch size type: %T", i)
	}

	if size == 0 {
		return errorsmod.Wrap(ErrInvalidBatchSize, "zero batch size")
	}

	return nil
}

func validateBatchIndex(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid batch index type: %T", i)
	}

	return nil
}

func validateVestingPeriod(i interface{}) error {
	period, ok := i.(int64)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "invalid vesting period type: %T", i)
	}

	if period == 0 {
		return errorsmod.Wrap(ErrInvalidVestingPeriod, "zero vesting period")
	}

	return nil
}
