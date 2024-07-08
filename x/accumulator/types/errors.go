package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/accumulator module sentinel errors
var (
	ErrInvalidPool     = sdkerrors.Register(ModuleName, 1100, "invalid pool address")
	ErrInvalidReceiver = sdkerrors.Register(ModuleName, 1101, "invalid receiver address")
)
