<!--
order: 1
-->

# Concepts

## The NFT module

The NFT module is primarily created for a special type of system of 
non-fungible tokens that hold locked native tokens and can be staked.


## Module requirements: 
* The module implements the non-fungible token based methods;
* Module configuration implies the tokens distribution in the genesis
block as well as token mint by the System master admin in the future;
* From the architectural perspective - each token represents an independent balance with locked tokens and tokens available for withdrawal;
* Besides owner and metadata URL, each token is characterized by a fixed `token_amount` (100%) and the unlock percent `unlock% = unlocked_amount / token_amount`. Every vesting period the `EndBlocker` unlocks a fixed step in tokens (`reward_per_period`); the last period unlocks only the remainder, so the unlock ends exactly at 100%;
* Each token can be delegated or undelegated using built-in methods, which means the standard delegation or undelegation of the whole token balance by the token owner;
* Each delegated token is able to collect rewards that can be immediately withdrawn by the token owner;
* If the token is not delegated, the owner can withdraw up to `token_amount × unlock% − withdrawn` until the token is fully unlocked. The limit is calculated from `token_amount`, not from the live balance, so rewards and donations do not increase it;
* Once the token is 100% unlocked, the owner can withdraw any part of the whole NFT balance, including rewards and donations;
* Delegated tokens can not be withdrawn: undelegate first (and wait for the unbonding to complete);
* Undelegations, staking rewards and regular bank transfers to the NFT address are donations. They do not change `token_amount`, the unlock percent or the vesting time;
* The vesting is topped up only with `MsgTopUp`, which can be sent by the NFT owner only. It increases `token_amount`; the unlock step and the already unlocked amount stay the same, so the vesting time grows proportionally to the top-up. A top-up of a fully unlocked NFT is rejected;
* The module admin can close an NFT with `MsgClose`: the whole NFT balance is sent to the owner without waiting for the end of vesting and the NFT is deleted. An NFT with delegations, unbonding delegations or redelegations can not be closed;
* If the token is delegated, it should not be able to transfer the token ownership;
* Token withdrawal implies the transfer of the corresponding amount of tokens from the NFT shard treasury to the token owner account;
* The locked amount in the minted NFTs can not exceed the amount in the NFT shard treasury;
* There should not be any other way to access the NFT shard treasury tokens besides withdrawing unlocked tokens from NFT.