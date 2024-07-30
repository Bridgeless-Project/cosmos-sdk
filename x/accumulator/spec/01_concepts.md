<!--
order: 1
-->

# Concepts

## The accumulator module
The accumulator module is designed to be the primary accumulator and
distributor of the native tokens in the system. According to System
tokenomics, all native tokens issued in the system will be minted to the 
accumulator module balance in the genesis block. By operating the
module configuration and business logic, the module will be able to 
dispose of funds according to System tokenomics.


Token distribution requirements:
* All tokens should be issued to the accumulator module account in the
genesis block;
* Module should be configured to split all tokens into three treasures:
  * community (admin) treasury; 
  * NFT shard treasury;
  * validator treasury;
* Validator treasury will be used to cover validators' rewards with linear distribution during 48 months according to the voting power;
* NFT shard treasury will be used by the NFT module and implies the locking tokens till some data in the non-fungible token can be staked;
* Admin treasury will have the self-inner distribution partially configured in the genesis block. This distribution implies the following distribution flows:
  * Tokens will be split into the launch and cliff parts;
  * Launch part of tokens should be transferred to the admin’s
  accounts in the genesis block;
  * Cliff parts of tokens will have different unlock periods (period
  after the distribution starts);
  * Cliff parts can optionally follow the linear distribution during
  the fixed period (by default - it should be the immediate
  transfer to the receiver);
  * Cliff parts receiver accounts can be optionally configured later
  with the System master admin;

Linear distribution or linear unlocking during a certain fixed period involves dividing the entire amount into several equal
parts and being able to transfer (or receive) these parts each time at the same interval.


 
 
 
 