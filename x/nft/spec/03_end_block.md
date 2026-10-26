<!--
order: 3
-->
# End Block: NFT Vesting Update

## Overview

The `EndBlocker` function is responsible for updating the vesting state of each NFT at the end of each block. This ensures that NFT owners can only withdraw the available amount as per the vesting schedule.

## NFT Vesting Information

Each NFT contains vesting information structured as follows:

- **Token Amount**: the fixed vesting entitlement (100%), changed only by `MsgTopUp`.
- **Unlocked Amount**: the cumulative unlocked amount (`unlock% = unlocked / token amount`).
- **Reward Per Period**: the fixed unlock step in tokens.
- **Last Vesting Block**: the block of the last vesting event.
- **Start Vesting Block**: the start block of the vesting.

Balance changes of the NFT address (rewards, undelegations, bank transfers) do not affect the vesting.

## Purpose of EndBlocker

The `EndBlocker` unlocks tokens of each NFT in batches according to the vesting schedule.

### Key Steps

1. **Retrieve NFTs**: The function retrieves the current batch of NFTs.
2. **Skip fully unlocked NFTs**: if `unlocked_amount >= token_amount`, the NFT is skipped.
3. **Check Start and Period**: the NFT is skipped until the start vesting block is reached and a full vesting period has passed since the last vesting block.
4. **Unlock**: it unlocks `min(reward_per_period × passed_periods, token_amount − unlocked_amount)`, so the last period unlocks only the remainder and the unlock ends exactly at 100%.
5. **Update Last Vesting Block**: it sets the last vesting block to the current one.

## Withdrawal limit

- Before the full unlock the owner can withdraw `min(unlocked_amount − withdrawn_amount, balance)`.
- After the full unlock the owner can withdraw any part of the whole NFT balance, including rewards and donations.
- A delegated NFT can not be withdrawn.
