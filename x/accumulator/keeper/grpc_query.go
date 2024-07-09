package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

var _ types.QueryServer = BaseKeeper{}
