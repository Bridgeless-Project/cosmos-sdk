package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// initialize validator distribution record
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}

// AfterValidatorRemoved performs clean up after a validator is removed
func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

// increment period
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

// create new delegation period record
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

// record the slash event
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error { return nil }
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationUpdated(ctx sdk.Context, delAddress sdk.AccAddress) error {
	nft, found := h.k.nftKeeper.GetNFT(ctx, delAddress.String())
	if found {
		delAddress = sdk.MustAccAddressFromBech32(nft.Owner)
	}

	proposals := h.k.GetProposalsFiltered(ctx, govtypes.QueryProposalsParams{
		Voter:          delAddress,
		ProposalStatus: govtypes.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD,
		Page:           1,
	})

	if len(proposals) > 0 {
		return stakingtypes.ErrVotedDelegation
	}

	return nil
}
