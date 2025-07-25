package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// TODO: Break into several smaller functions for clarity

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, proposal v1.Proposal) (passes bool, burnDeposits bool, tallyResults v1.TallyResult) {
	results := make(map[v1.VoteOption]sdk.Dec)
	results[v1.OptionYes] = sdk.ZeroDec()
	results[v1.OptionAbstain] = sdk.ZeroDec()
	results[v1.OptionNo] = sdk.ZeroDec()
	results[v1.OptionNoWithVeto] = sdk.ZeroDec()

	tallyParams := keeper.GetTallyParams(ctx)

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]v1.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	keeper.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = v1.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			v1.WeightedVoteOptions{},
		)

		return false
	})

	keeper.IterateVotes(ctx, proposal.Id, func(vote v1.Vote) bool {
		// if validator, just record it in the map
		voter := sdk.MustAccAddressFromBech32(vote.Voter)

		valAddrStr := sdk.ValAddress(voter.Bytes()).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		// get all nfts related for delegator
		nfts, _, err := keeper.nftKeeper.GetAllNFTsByOwnerWithPagination(ctx, voter.String(), &query.PageRequest{Limit: query.MaxLimit})
		if err != nil {
			keeper.Logger(ctx).Error("failed to get all nfts for the validator ", err)
			return false
		}

		// iterate over all delegator's nfts to calculate power
		for _, nft := range nfts {
			// iterate over all delegations from voter, deduct from any delegated-to validators
			keeper.sk.IterateDelegations(ctx, sdk.MustAccAddressFromBech32(nft.Address), func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
				valAddrStr = delegation.GetValidatorAddr().String()

				// validate that delegation is available to vote
				if !delegation.GetTimestamp().Add(tallyParams.DepositLockingPeriod).Before(ctx.BlockTime()) {
					keeper.Logger(ctx).Info(fmt.Sprintf("delegation %s is not yet unlocked", delegation.GetValidatorAddr().String()))
					return false
				}

				if val, ok := currValidators[valAddrStr]; ok {
					// There is no need to handle the special case that validator address equal to voter address.
					// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
					val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
					currValidators[valAddrStr] = val

					// delegation shares * bonded / total shares
					votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

					for _, option := range vote.Options {
						weight, _ := sdk.NewDecFromStr(option.Weight)
						subPower := votingPower.Mul(weight)
						results[option.Option] = results[option.Option].Add(subPower)
					}
					totalVotingPower = totalVotingPower.Add(votingPower)
				}

				return false
			})
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		keeper.sk.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr = delegation.GetValidatorAddr().String()

			// validate that delegation is available to vote
			if !delegation.GetTimestamp().Add(tallyParams.DepositLockingPeriod).Before(ctx.BlockTime()) {
				keeper.Logger(ctx).Info(fmt.Sprintf("delegation %s is not yet unlocked", delegation.GetValidatorAddr().String()))
				return false
			}

			if val, ok := currValidators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

				// delegation shares * bonded / total shares
				votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, option := range vote.Options {
					weight, _ := sdk.NewDecFromStr(option.Weight)
					subPower := votingPower.Mul(weight)
					results[option.Option] = results[option.Option].Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})

		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			weight, _ := sdk.NewDecFromStr(option.Weight)
			subPower := votingPower.Mul(weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyResults = v1.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if keeper.sk.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(sdk.NewDecFromInt(keeper.sk.TotalBondedTokens(ctx)))
	quorum, _ := sdk.NewDecFromStr(tallyParams.Quorum)
	if percentVoting.LT(quorum) {
		return false, false, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[v1.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	vetoThreshold, _ := sdk.NewDecFromStr(tallyParams.VetoThreshold)
	if results[v1.OptionNoWithVeto].Quo(totalVotingPower).GT(vetoThreshold) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	threshold, _ := sdk.NewDecFromStr(tallyParams.Threshold)
	if results[v1.OptionYes].Quo(totalVotingPower.Sub(results[v1.OptionAbstain])).GT(threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}
