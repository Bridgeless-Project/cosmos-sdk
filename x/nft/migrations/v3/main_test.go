package v3

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v3types "github.com/cosmos/cosmos-sdk/x/nft/migrations/v3/types"
)

func coin(n int64) sdk.Coin { return sdk.NewInt64Coin("stake", n) }

func TestConvert(t *testing.T) {
	cases := []struct {
		name                                       string
		count, reward, token, available            int64
		wantUnlocked, wantWithdrawn, wantAvailable int64
	}{
		{"nothing withdrawn", 3, 10, 100, 30, 30, 0, 30},
		{"partly withdrawn", 3, 10, 100, 12, 30, 18, 12},
		{"fully unlocked, remainder period", 11, 10, 105, 105, 105, 0, 105},
		{"count capped by token", 20, 10, 100, 40, 100, 60, 40},
		{"inconsistent available", 1, 10, 100, 25, 25, 0, 25},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convert(v3types.NFT{
				Address:             "addr",
				Denom:               "stake",
				VestingPeriodsCount: tc.count,
				RewardPerPeriod:     coin(tc.reward),
				TokenAmount:         coin(tc.token),
				AvailableToWithdraw: coin(tc.available),
			})

			require.Equal(t, tc.wantUnlocked, got.UnlockedAmount.Amount.Int64())
			require.Equal(t, tc.wantWithdrawn, got.WithdrawnAmount.Amount.Int64())
			require.Equal(t, tc.wantAvailable, got.WithdrawableBeforeUnlock().Int64())
		})
	}
}
