package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	"github.com/stretchr/testify/require"
)

func newNFT(token, step, unlocked, withdrawn int64) types.NFT {
	coin := func(v int64) sdk.Coin { return sdk.NewInt64Coin("stake", v) }

	return types.NFT{
		Denom:           "stake",
		TokenAmount:     coin(token),
		RewardPerPeriod: coin(step),
		UnlockedAmount:  coin(unlocked),
		WithdrawnAmount: coin(withdrawn),
	}
}

func TestUnlockStepEndsExactlyAtFull(t *testing.T) {
	nft := newNFT(100, 30, 0, 0)

	for i := 0; i < 3; i++ {
		require.Equal(t, sdk.NewInt(30), nft.UnlockStep())
		nft.UnlockedAmount = nft.UnlockedAmount.AddAmount(nft.UnlockStep())
	}

	require.False(t, nft.IsFullyUnlocked())
	require.Equal(t, sdk.NewInt(10), nft.UnlockStep())

	nft.UnlockedAmount = nft.UnlockedAmount.AddAmount(nft.UnlockStep())
	require.True(t, nft.IsFullyUnlocked())
	require.Equal(t, sdk.NewInt(100), nft.UnlockedAmount.Amount)
	require.True(t, nft.UnlockStep().IsZero())
}

func TestWithdrawableBeforeUnlock(t *testing.T) {
	require.Equal(t, sdk.NewInt(40), newNFT(100, 20, 60, 20).WithdrawableBeforeUnlock())
	require.True(t, newNFT(100, 20, 20, 30).WithdrawableBeforeUnlock().IsZero())
}

func TestTopUpKeepsStepAndExtendsVesting(t *testing.T) {
	nft := newNFT(100, 25, 50, 0)
	nft.TokenAmount = nft.TokenAmount.AddAmount(sdk.NewInt(50))

	require.False(t, nft.IsFullyUnlocked())
	require.Equal(t, sdk.NewInt(25), nft.RewardPerPeriod.Amount)
	require.Equal(t, sdk.NewInt(50), nft.UnlockedAmount.Amount)
}

func TestMsgValidateBasic(t *testing.T) {
	addr := sdk.AccAddress("addr_______________").String()

	require.NoError(t, types.NewMsgClose(addr, addr).ValidateBasic())
	require.Error(t, types.NewMsgClose("bad", addr).ValidateBasic())

	require.NoError(t, types.NewMsgTopUp(addr, addr, sdk.NewInt64Coin("stake", 1)).ValidateBasic())
	require.Error(t, types.NewMsgTopUp(addr, addr, sdk.NewInt64Coin("stake", 0)).ValidateBasic())
}
