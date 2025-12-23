package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ banktypes.BankHooks = Hooks{}

// Hooks creates new nft hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// BeforeSendTokenToAddress handles cases when recipient is a NFT
func (h Hooks) BeforeSendTokenToAddress(_ sdk.Context, _, _ sdk.Address, _ sdk.Coins) error {
	// this hook is ignored
	return nil
}

// AfterSendTokenToAddress updates vesting params for NFT
func (h Hooks) AfterSendTokenToAddress(ctx sdk.Context, receiver sdk.Address, amt sdk.Coins) error {
	params := h.k.GetParams(ctx)

	nft, found := h.k.GetNFT(ctx, receiver.String())
	if !found {
		return nil
	}

	// get count of additional periods
	amount := amt.AmountOf(nft.Denom).BigInt()
	reward := nft.RewardPerPeriod.Amount.BigInt()

	full := new(big.Int)
	temp := new(big.Int).Add(amount, reward)
	temp.Sub(temp, big.NewInt(1))
	full.Quo(temp, reward)

	nft.VestingPeriodsLimit += full.Int64()

	nft.TokenAmount = nft.TokenAmount.Add(sdk.NewCoin(nft.Denom, amt.AmountOf(nft.Denom)))

	blockDistanceFromLastVesting := ctx.BlockHeight() - nft.LastVestingBlock

	// timeToRestartVesting includes time to get to batch where nft will be processed and covers the time was spent after last vesting block
	timeToRestartVesting := h.k.GetNftBatchBlockDistance(ctx, nft.Address) + uint64(blockDistanceFromLastVesting)

	// adding the time of vesting
	nft.TotalVestingTime += full.Int64()*params.VestingPeriod + int64(timeToRestartVesting)

	h.k.SetNFT(ctx, nft)

	return nil
}
