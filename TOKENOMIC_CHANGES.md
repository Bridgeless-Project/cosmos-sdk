# Tokenomic

## Overview

- Limit tokens amount 1000000000 BRIDGE and additionally each token can be split to 10^18 abridge tokens (BRIDGE has 18
  decimals)
- All tokens are split to 3 pools and managed by `accumulator` module
- Rewards by a block can be taken only from the `accumulator` module
- System operates a new `NFT` module. The `NFT` module holds locked native tokens that can be staked, unstaked, etc.
- Staked amount from a NFT have the same power as the same amount of native tokens. But delegator will get +20% (
  configured) by stake from NFT
- Each NFT has a vesting time to unlock balance linearly to withdraw

## Changes in default modules

### Proto files

The `distribution` module params proto file:

    message Params {
        ...
        string nft_proposer_reward = 4 [
            (cosmos_proto.scalar)  = "cosmos.Dec",
            (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
            (gogoproto.nullable)   = false
        ];
        ...

Define a new multiplier for rewards by NFT's staking. In our case this field is equal to 0.2000000...(20%)

One more change related to proto file: `staking` module params:

    message Delegation {
        ...
        string amount = 4 [
            (cosmos_proto.scalar)  = "cosmos.Dec",
            (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
            (gogoproto.nullable)   = false
        ];
    }

The additional field is used to store count of staked tokens in raw representation (__not shares__) (count of tokens
sent from delegator
to validator).

The `mint` module proto:

    message Params {
        option (gogoproto.goproto_stringer) = false;
        
        // type of coin to mint
        string mint_denom = 1;
        // expected blocks per month
        uint64 blocks_per_month = 2;
        // block when no additional tokens will be minted
        uint64 end_block = 3;
        
        cosmos.base.v1beta1.Coin month_reward = 4
        [(gogoproto.nullable) = false, (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Coin"];
    }

The additional fields describe custom block rewards policy (linear distribution with constant amount per block).

## Updates of business logic

### Distribution module

- [calculateDelegationRewardsBetween](x/distribution/keeper/delegation.go)
  This function calculates the rewards accrued by a delegation between two periods.

        func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI,
        endingPeriod uint64) (rewards sdk.DecCoins) {

        ...
    
        params := k.GetParams(ctx)
        _, isNFTstake := k.nftKeeper.GetNFT(ctx, del.GetDelegatorAddr().String()) // Checks if this delegation is NFT
    
        wrapStake := func(stake sdk.Dec) sdk.Dec { // Special wrapper to add additional 20% profit by NFT staking.
            if isNFTstake {
                return stake.Add(stake.Mul(params.NftProposerReward))
            }
            return stake
        }
    
        ...
        
        endingHeight := uint64(ctx.BlockHeight())
        if endingHeight > startingHeight {
            k.IterateValidatorSlashEventsBetween(ctx, del.GetValidatorAddr(), startingHeight, endingHeight,
                func(height uint64, event types.ValidatorSlashEvent) (stop bool) {
                    if event.ValidatorPeriod > startingPeriod {
                        // Using wrapped stake amount
                        rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, wrapStake(stake))...) 
                        
                        ...
                    }
                    return false
                },
            )
        }
    
        ....
    
        // Calculate rewards for final period (also using wrapped stake amount)
        rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, wrapStake(stake))...)
        return rewards

}

The key update in this function is the anonymous function `wrapStake(sdk.Dec)`. The `wrapStake` function checks whether
the delegation is made by an NFT or not, and based on this, it updates the rewards accordingly.

- [IncrementValidatorPeriod](x/distribution/keeper/validator.go)
  Increment validator period, returning the period just ended

Since the total amount of rewards remains constant, but NFT delegators are entitled to receive a higher proportion of
tokens, it is necessary to normalize the rewards.

    func (k Keeper) IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64 {
       
        ...
    
        params := k.GetParams(ctx)
        // Iterating over all delegations and adding additional 20% profit for NFT delegations to take into account in the total stake amount divisor.
        for _, del := range k.stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator()) {
            if _, found := k.nftKeeper.GetNFT(ctx, del.GetDelegatorAddr().String()); found {
                tokens = tokens.Add(del.Amount.Add(del.Amount.Mul(params.NftProposerReward)))
                continue
            }
    
            tokens = tokens.Add(del.Amount)
        }
    
        if tokens.IsZero() {

            ...

        } else {
            current = rewards.Rewards.QuoDecTruncate(tokens)
        }
    
        ...
    }

The main update in the function is the following code segment, where all rewards are summed up in relation to
the delegator's multiplier. All operations are adjusted based on the raw staked token amounts.

    params := k.GetParams(ctx)
    for _, del := range k.stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator()) {
        if _, found := k.nftKeeper.GetNFT(ctx, del.GetDelegatorAddr().String()); found {
            tokens = tokens.Add(del.Amount.Add(del.Amount.Mul(params.NftProposerReward)))
            continue
        }

        tokens = tokens.Add(del.Amount)
    }

### Staking module

- [Delegation](proto/cosmos/staking/v1beta/staking.proto)
  Expand the `Delegation` struct by adding a timestamp field. The field includes the timestamp when the deposit was updated.

      google.protobuf.Timestamp timestamp = 5 [(gogoproto.stdtime) = true, (gogoproto.nullable)   = false] ;

- [BeforeDelegationUpdated](x/staking/types/expected_keepers.go) This hook, generally, is used to validate that delegation can be updated.


- [Delegate](x/staking/keeper/delegation.go)
  The `Delegate` performs a delegation, set/update everything necessary within the store. The `tokenSrc` indicates the
  bond status of the incoming funds.

      func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc types.BondStatus,
          validator types.Validator, subtractAccount bool) (newShares sdk.Dec, err error) {

        ...

        delegation, found := k.GetDelegation(ctx, delAddr, validator.GetOperator())
        if !found {
            delegation = types.NewDelegation(delAddr, validator.GetOperator(), sdk.ZeroDec(), sdk.ZeroDec())
        }
     
        ...

        // Update delegation
        delegation.Shares = delegation.Shares.Add(newShares)
        delegation.Amount = delegation.Amount.Add(sdk.NewDecFromInt(bondAmt))
        k.SetDelegation(ctx, delegation)
    
       ...
  }

Set the delegate amount. In the first case, if no delegation is found, initialize the amount with an empty value. Just
before returning, set `delegation.Amount`.


- [Unbond](x/staking/keeper/delegation.go)
  The `Unbond` method unbonds a particular delegation and perform associated store operations.

        func (k Keeper) Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (amount math.Int, err error) {

          ...
      
          validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)
          if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
              // if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
              k.RemoveValidator(ctx, validator.GetOperator())
          }
      
          // Using the amount received from RemoveValidatorTokensAndShares before
          delegation.Amount = delegation.Amount.Sub(sdk.NewDecFromInt(amount))


          if delegation.Shares.IsZero() {
              err = k.RemoveDelegation(ctx, delegation)
          } else {
              k.SetDelegation(ctx, delegation)
              // call the after delegation modification hook
              err = k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
          }
      
          if err != nil {
              return amount, err
          }
      
          return amount, nil
        }

Move the call to `RemoveValidatorTokensAndShares` earlier in the process. This is necessary to set the amount of tokens
before updating the delegation.

- [Undelegate](x/staking/keeper/delegation.go) Insert `BeforeDelegationUpdated` hook to validate that delegation is not used for voting

      func (k Keeper) Undelegate(
        ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
      ) (time.Time, error) {
        if err := k.hooks.BeforeDelegationUpdated(ctx, delAddr); err != nil {
          return time.Time{}, err
        }
      
        ...
 
      }


- [BeginRedelegation](x/staking/keeper/delegation.go) Insert `BeforeDelegationUpdated` hook to validate that delegation is not used for voting

      func (k Keeper) BeginRedelegation(
        ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
      ) (completionTime time.Time, err error) {
        
          ...
      
          if err := k.hooks.BeforeDelegationUpdated(ctx, delAddr); err != nil {
            return time.Time{}, err
          }
          ...
      }


### Mint module

- [SendFromAccumulator](x/mint/keeper/keeper.go) This function is used to sent tokens from validator pool to the
  accumulator module.

        func (k Keeper) SendFromAccumulator(ctx sdk.Context, amount sdk.Coins) error {
          err := k.accumulatorKeeper.DistributeToModule(ctx, accumulatortypes.ValidatorPoolName, amount, types.ModuleName)
          if err != nil {
            err = errors.Wrap(err, "failed to call accumulator module")
            k.Logger(sdk.UnwrapSDKContext(ctx)).Error(err.Error())
            return err
          }
      
          return nil
        }

This approach ensures that no new tokens will be created, and the amount of rewards for the validator
is properly managed by the validator pool.

- [BeginBlocker](x/mint/abci.go)
  The `BeginBlocker` mints new tokens for the previous block with custom policy.

      func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

          defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
          params := k.GetParams(ctx)

            // skip if all tokens already minted
            if uint64(ctx.BlockHeight()) >= params.EndBlock {
                return
            }
            monthReward := sdk.NewDecFromInt(params.MonthReward.Amount)
            mintedAmount := monthReward.QuoInt(sdk.NewInt(int64(params.BlocksPerMonth)))
        
            // mint coins, update supply
            mintedCoin := sdk.NewCoin(params.MintDenom, mintedAmount.TruncateInt())
            mintedCoins := sdk.NewCoins(mintedCoin)
        
            err := k.SendFromAccumulator(ctx, mintedCoins)
            if err != nil {
                k.Logger(ctx).Error("failed to send tokens from accumulator")
                return
            }
        
            // send the minted coins to the fee collector account
            err = k.AddCollectedFees(ctx, mintedCoins)
            if err != nil {
                panic(err)
            }
        
            if mintedCoin.Amount.IsInt64() {
                defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
            }
        
            ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                    types.EventTypeMint,
                    sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
                ),
            )
      }

The following code snippet illustrates the process of minting tokens.


### Gov module 

- [Tally](x/gov/keeper/tally.go) Tally iterates over the votes and updates the tally of a proposal based on the voting power of the voters

      func (keeper Keeper) Tally(ctx sdk.Context, proposal v1.Proposal) (passes bool, burnDeposits bool, tallyResults v1.TallyResult) {
          
          ...
      
          votingParams := keeper.GetVotingParams(ctx)
      
          ...
      
          keeper.IterateVotes(ctx, proposal.Id, func(vote v1.Vote) bool {
              ...
      
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
      
                      // validate that delegation is available to votes
                      if !delegation.GetTimestamp().Add(votingParams.LockingPeriod).Before(ctx.BlockTime()) {
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
                  if !delegation.GetTimestamp().Add(votingParams.LockingPeriod).Before(ctx.BlockTime()) {
                      keeper.Logger(ctx).Info(fmt.Sprintf("delegation %s is not yet unlocked", delegation.GetValidatorAddr().String()))
                      return false
                  }
      
                  ...
              })
      
              keeper.deleteVote(ctx, vote.ProposalId, voter)
              return false
          })

  Add ability to vote taking into account NFT delegation. And also to prevent `sandwich attack` this function ignores recent delegations
