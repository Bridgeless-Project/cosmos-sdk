<!--
order: 3
-->

## Overview

The `EndBlocker` function handles the vesting process for administrators (admins) by updating their vesting state at the end of each block. 
This ensures that admins receive their allocated rewards based on the vesting schedule. Locked tokens are stored on accumulator pull address.

## Purpose of EndBlocker

The `EndBlocker` function updates the vesting state for each admin, distributing the rewards to their accounts based on the vesting schedule.

### Key Steps

1. **Retrieve Admins**: The function retrieves all admins.
2. **Check Vesting Period**: It checks if the current block time is greater than or equal to the last vesting time plus the vesting period.
3. **Check Vesting Counter**: It ensures that the vesting counter has not reached the maximum number of vesting periods.
4. **Distribute Rewards**: If both conditions are met, it distributes the reward to the admin's account.
5. **Increment Counter**: It increments the vesting counter and updates the last vesting time.


The `EndBlocker` function is crucial for maintaining the correct vesting state for admins, ensuring that they receive the amount that has
vested over time according to the specified schedule.