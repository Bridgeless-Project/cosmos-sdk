# Tokenomic

## Overview

-   Limit tokens amount 1000000000 BRIDGE and additionally each token can be split to  10^18 abridge tokens
-   System contains a NFT. The NFT holds locked native tokens and can be staked
-   Staked amount from a nft have the same power as the same amount of native tokens. But delegator will get +20% by stake from nft
-   Each nft has a vesting time to unlock balance to withdraw

## changes in default modules

### Proto files 

Distribution params proto file

    message Params {
        ...
        string nft_proposer_reward = 4 [
        (cosmos_proto.scalar)  = "cosmos.Dec",
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable)   = false
        ];
        ...
 
Define a new multiplier for rewards by nft staking.  In our case this field is equal to 0.2000000...(20%)

One more change related to proto file
Staking params proto file

    message Delegation {
        ...
        string amount = 4 [
        (cosmos_proto.scalar)  = "cosmos.Dec",
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable)   = false
        ];
    }

The additional field is used to store count of staked tokens in raw representation (count of tokens sent from  delegator to validator) 

### Updates of usiness logic 

Distribution module 

- [calculateDelegationRewardsBetween](x/distribution/keeper/delegation.go)
The function calculates the rewards accrued by a delegation between two periods



    // calculate the total rewards accrued by a delegation
    func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins) {

        ...
    
        params := k.GetParams(ctx)
        _, isNFTstake := k.nftKeeper.GetNFT(ctx, del.GetDelegatorAddr().String())
    
        wrapStake := func(stake sdk.Dec) sdk.Dec {
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
                        rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, wrapStake(stake))...)
                        
                        ...
                    }
                    return false
                },
            )
        }
    
        ....
    
        // calculate rewards for final period
        rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, wrapStake(stake))...)
        return rewards
    }

The key update in this function is the anonymous function `wrapStake(sdk.Dec)`. The `wrapStake` function checks whether the delegation
is made by an NFT or not, and based on this, it updates the rewards accordingly.




- [IncrementValidatorPeriod](x/distribution/keeper/validator.go)
Increment validator period, returning the period just ended

Since the total amount of rewards remains constant, but NFT delegators are entitled to receive a higher proportion of tokens, it is necessary to normalize the rewards.
  
    func (k Keeper) IncrementValidatorPeriod(ctx sdk.Context, val stakingtypes.ValidatorI) uint64 {
       
        ...
    
        params := k.GetParams(ctx)
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




Staking module 

- [Delegate](x/staking/keeper/delegation.go)
Delegate performs a delegation, set/update everything necessary within the store. tokenSrc indicates the bond status of the incoming funds.
    

    func (k Keeper) Delegate(
    ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc types.BondStatus,
    validator types.Validator, subtractAccount bool,
    ) (newShares sdk.Dec, err error) {
    
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

Set the delegate amount. In the first case, if no delegation is found, initialize the amount with an empty value. Just before returning, set `delegation.Amount`


- [Unbond](x/staking/keeper/delegation.go)

Unbond unbonds a particular delegation and perform associated store operations.
    

    func (k Keeper) Unbond(
    ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
    ) (amount math.Int, err error) {
        
        ...
    
        validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)
        if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
            // if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
            k.RemoveValidator(ctx, validator.GetOperator())
        }
    
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


Move the call to RemoveValidatorTokensAndShares earlier in the process. This is necessary to set the amount of tokens before updating the delegation.






