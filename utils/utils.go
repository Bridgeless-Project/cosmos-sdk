package utils

import (
	"github.com/cosmos/cosmos-sdk/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewNativeTokens(amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(types.NativeToken, sdk.NewInt(amount)))
}
