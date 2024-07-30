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
* Besides owner and metadata URL, each token is characterized by the number of locked tokens and unlock period during which tokens will be unlocked linearly;
* Each token can be delegated or undelegated using built-in methods, which means the standard delegation or undelegation of the whole token balance by the token owner;
* Each delegated token is able to collect rewards that can be immediately withdrawn by the token owner;
* If the token is not delegated and contains unlocked tokens - they can be immediately withdrawn from the token balance;
* If the token is delegated, it should not be able to transfer the token ownership;
* Token withdrawal implies the transfer of the corresponding amount of tokens from the NFT shard treasury to the token owner account;
* The locked amount in the minted NFTs can not exceed the amount in the NFT shard treasury;
* There should not be any other way to access the NFT shard treasury tokens besides withdrawing unlocked tokens from NFT.