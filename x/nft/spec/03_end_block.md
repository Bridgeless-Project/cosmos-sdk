<!--
order: 3
-->
# End Block: NFT Vesting Update

## Overview

The `EndBlocker` function is responsible for updating the vesting state of each NFT at the end of each block. This ensures that NFT owners can only withdraw the available amount as per the vesting schedule.

## NFT Vesting Information

Each NFT contains vesting information structured as follows:

- **Vesting Periods Count**: The current count of how many periods have vested.
- **Available to Withdraw**: The amount available for withdrawal.
- **Last Vesting Time**: The block of the last vesting event.
- **Start Vesting Time** The start block of the vesting.

## Purpose of EndBlocker

The `EndBlocker` function updates the vesting state for each NFT, ensuring that the available amount to withdraw is updated based on the vesting schedule.

### Key Steps

1. **Retrieve NFTs**: The function retrieves all NFTs.
2. **Check Vesting Period**: It checks if the current block time is greater than or equal to the last vesting time plus the vesting period.
3. **Check Vesting Counter**: It ensures that the vesting counter has not reached the maximum number of vesting periods.
4. **Update Available Amount**: If both conditions are met, it calculates the new available amount to withdraw based on the reward per period and updates the `Available to Withdraw` field.
5. **Increment Counter**: It increments the vesting counter and updates the last vesting time.

