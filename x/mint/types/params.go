package types

import (
	"errors"
	"fmt"
	"sigs.k8s.io/yaml"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultHalvingBlocks     = 1000000
	DefaultMaxHalvingPeriods = 7
)

// Parameter store keys
var (
	KeyMintDenom            = []byte("MintDenom")
	KeyHalvingBlocks        = []byte("HalvingBlocks")
	KeyMaxHalvingPeriods    = []byte("MaxHalvingPeriods")
	KeyCurrentHalvingPeriod = []byte("CurrentHalvingPeriod")
	KeyBlockReward          = []byte("BlockReward")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, halvingBlocks uint64, blockReward sdk.Coin, currentHalvingPeriod, maxHalvingPeriods uint32,
) Params {
	return Params{
		MintDenom:            mintDenom,
		HalvingBlocks:        halvingBlocks,
		BlockReward:          blockReward,
		CurrentHalvingPeriod: currentHalvingPeriod,
		MaxHalvingPeriods:    maxHalvingPeriods,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:            sdk.DefaultBondDenom,
		HalvingBlocks:        DefaultHalvingBlocks,
		BlockReward:          sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(6)),
		CurrentHalvingPeriod: 0,
		MaxHalvingPeriods:    DefaultMaxHalvingPeriods,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}

	if err := validateHalvingBlocks(p.HalvingBlocks); err != nil {
		return err
	}

	if err := validateCurrentHalvingPeriod(p.CurrentHalvingPeriod); err != nil {
		return err
	}

	if err := validateBlockReward(p.BlockReward); err != nil {
		return err
	}

	if err := validateMaxHalvingPeriods(p.MaxHalvingPeriods); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyHalvingBlocks, &p.HalvingBlocks, validateHalvingBlocks),
		paramtypes.NewParamSetPair(KeyMaxHalvingPeriods, &p.MaxHalvingPeriods, validateMaxHalvingPeriods),
		paramtypes.NewParamSetPair(KeyCurrentHalvingPeriod, &p.CurrentHalvingPeriod, validateCurrentHalvingPeriod),
		paramtypes.NewParamSetPair(KeyBlockReward, &p.BlockReward, validateBlockReward),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateHalvingBlocks(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return errors.New("halving blocks cannot be 0")
	}
	return nil
}

func validateMaxHalvingPeriods(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return errors.New("max halving periods cannot be 0")
	}
	return nil
}

func validateCurrentHalvingPeriod(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return errors.New(fmt.Sprintf("invalid parameter type: %T", i))
	}

	return nil
}

func validateBlockReward(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return errors.New(fmt.Sprintf("invalid parameter type: %T", i))
	}

	if v.IsNegative() {
		return errors.New("block reward cannot be negative")
	}

	if !v.IsValid() {
		return errors.New(fmt.Sprintf("invalid block reward: %s", v))
	}

	return nil
}
