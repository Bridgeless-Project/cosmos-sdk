package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestVestingParamsValidation(t *testing.T) {
	for _, v := range []int64{0, -1} {
		require.Error(t, validateVestingTime(v))
		require.Error(t, validateVestingPeriod(v))
	}
	require.NoError(t, validateVestingTime(int64(10)))
	require.NoError(t, validateVestingPeriod(int64(3)))

	p := DefaultParams()
	p.ModuleAdmin = sdk.AccAddress("admin_______________").String()
	p.TotalVestingTime, p.VestingPeriod = 10, 3
	require.NoError(t, p.Validate())
	p.TotalVestingTime, p.VestingPeriod = 2, 3
	require.Error(t, p.Validate())
}
