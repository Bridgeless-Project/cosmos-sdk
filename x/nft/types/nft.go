package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// IsFullyUnlocked returns true once the whole token amount is unlocked.
func (n NFT) IsFullyUnlocked() bool {
	return n.UnlockedAmount.Amount.GTE(n.TokenAmount.Amount)
}

// UnlockStep returns the amount to unlock in one period: RewardPerPeriod, but not more than the locked remainder.
func (n NFT) UnlockStep() sdk.Int {
	remainder := n.TokenAmount.Amount.Sub(n.UnlockedAmount.Amount)
	if remainder.IsNegative() {
		return sdk.ZeroInt()
	}

	return sdk.MinInt(n.RewardPerPeriod.Amount, remainder)
}

// WithdrawableBeforeUnlock returns TokenAmount × unlock% − withdrawn, i.e. unlocked − withdrawn.
func (n NFT) WithdrawableBeforeUnlock() sdk.Int {
	amount := n.UnlockedAmount.Amount.Sub(n.WithdrawnAmount.Amount)
	if amount.IsNegative() {
		return sdk.ZeroInt()
	}

	return amount
}
