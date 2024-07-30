
# Client

## Query


### CLI

A user can query and interact with the `nft` module using the CLI.

The `query` commands allow users to query `nft` state.

```sh
simd query nft --help
```

#### nfts

The `nfts` command allow users to query all nfts

```sh
simd query nft nfts [flags]
```

Example:

```sh
simd query nft nfts
```

Example Output:

```sh
nft:
- address: bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37
  available_to_withdraw:
    amount: "0"
    denom: abridge
  denom: abridge
  last_vesting_time: "1722244515"
  owner: bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4
  reward_per_period:
    amount: "800000000000000000000"
    denom: abridge
  uri: ""
  vesting_counter: "2"
  vesting_period: "86400"
  vesting_periods_count: "100"
- address: bridge100rx5aqsunw3ta9nclsakkkj3t7sl5dzrza2gsw2n6hf62r2sswqvwcxa4
  available_to_withdraw:
    amount: "0"
    denom: abridge
  denom: abridge
  last_vesting_time: "1722244515"
  owner: bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4
  reward_per_period:
    amount: "800000000000000000000"
    denom: abridge
  uri: ""
  vesting_counter: "2"
  vesting_period: "86400"
  vesting_periods_count: "100"
```

#### nft 

The `nft` command allow users to query the nft

```sh
simd query nft nft [flags]
```

Example:

```sh
simd query nft nft bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37
```

Example Output:

```sh
- address: bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37
  available_to_withdraw:
    amount: "0"
    denom: abridge
  denom: abridge
  last_vesting_time: "1722244515"
  owner: bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4
  reward_per_period:
    amount: "800000000000000000000"
    denom: abridge
  uri: ""
  vesting_counter: "2"
  vesting_period: "86400"
  vesting_periods_count: "100"
```



#### owners 

The `owners` command allow users to query the addresses of nft holdres

```sh
simd query nft owners [flags]
```

Example:

```sh
simd query nft owners 
```

Example Output:

```sh
- owners: bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37
```

#### owner

The `owner` command allow users to query the nfts by holder address

```sh
simd query nft owners [flags]
```

Example:

```sh
simd query nft owners 
```

Example Output:

```sh
- owners: bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37
```


## HTTP

A user can query and interact with the `nft` module using the HTTP.

#### nfts

The `nfts` endpoint allow users to query all nfts

```
/cosmos/nfts
```

Example:

```
curl -X GET "https://rpc-api.node1.devnet.bridgeless.com/cosmos/nfts" -H  "accept: application/json"
```

Example Output:

```
{
  "nft": [
    {
      "address": "bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37",
      "owner": "bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4",
      "uri": "",
      "vesting_period": "86400",
      "reward_per_period": {
        "denom": "abridge",
        "amount": "800000000000000000000"
      },
      "vesting_periods_count": "100",
      "available_to_withdraw": {
        "denom": "abridge",
        "amount": "0"
      },
      "last_vesting_time": "1722262887",
      "vesting_counter": "2",
      "denom": "abridge"
    }
  ]
}
```

#### nft 

The `nft` endpoint allow users to query the nft

```
/cosmos/nfts/{address}
```

Example:

```
curl -X GET "https://rpc-api.node1.devnet.bridgeless.com/cosmos/nfts/bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37" -H  "accept: application/json"
```

Example Output:

```
{
  "nft": {
    "address": "bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37",
    "owner": "bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4",
    "uri": "",
    "vesting_period": "86400",
    "reward_per_period": {
      "denom": "abridge",
      "amount": "800000000000000000000"
    },
    "vesting_periods_count": "100",
    "available_to_withdraw": {
      "denom": "abridge",
      "amount": "0"
    },
    "last_vesting_time": "1722262887",
    "vesting_counter": "2",
    "denom": "abridge"
  }
}
```



#### owners 

The `owners` endpoint allow users to query the addresses of nft holdres

```
/cosmos/owners
```

Example:

```
curl -X GET "https://rpc-api.node1.devnet.bridgeless.com/cosmos/owners" -H  "accept: application/json"
```

Example Output:

```
{
  "owner": [
    "bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4",
    "bridge10awan8e059lkt7ejyu6qlkwwngs6yxxe4x5lhslpdvqlcmy3nrasd6zl5f"
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

#### owner

The `owner` endpoint allow users to query the nfts by holder address

```sh
/cosmos/nfts/{address}
```

Example:

```sh
curl -X GET "https://rpc-api.node1.devnet.bridgeless.com/cosmos/nfts/bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37" -H  "accept: application/json"
```

Example Output:

```
{
  "nft": {
    "address": "bridge100mehmdwp8kwlj7aak450cu02u9d32cza95yrv8q776sytjfg3mstssj37",
    "owner": "bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4",
    "uri": "",
    "vesting_period": "86400",
    "reward_per_period": {
      "denom": "abridge",
      "amount": "800000000000000000000"
    },
    "vesting_periods_count": "100",
    "available_to_withdraw": {
      "denom": "abridge",
      "amount": "0"
    },
    "last_vesting_time": "1722262887",
    "vesting_counter": "2",
    "denom": "abridge"
  }
}
```


## Message

```protobuf
service Msg {
  rpc Send(MsgSend) returns (MsgSendResponse);

  rpc Withdraw (MsgWithdrawal) returns (MsgWithdrawalResponse);

  rpc Undelegate(MsgUndelegate) returns (MsgUndelegateResponse);

  rpc Delegate(MsgDelegate) returns (MsgDelegateResponse);
  
  rpc Redelegate ( MsgRedelegate) returns (MsgRedelegateResponse);
}
```

#### Send

This message is used to change a nft owner
```protobuf
message MsgSend {
  string creator = 1;
  string address = 2;
  string recipient = 3;
  cosmos.base.v1beta1.Coin amount = 5 [(gogoproto.nullable) = false];
}

message MsgSendResponse {}
```

#### Withdrawal

This message is used to withdraw available(unlocked) nft amount 

```protobuf
message MsgWithdrawal {
  string creator = 1;
  string address = 2;
  cosmos.base.v1beta1.Coin amount = 5 [(gogoproto.nullable) = false];
}

message MsgWithdrawalResponse {}
```

#### Delegate

This message is used to stake tokens from nft balance to a validator 

```protobuf
message MsgDelegate {
  string creator = 1;
  string address = 2;
  string validator = 3;
  cosmos.base.v1beta1.Coin amount = 5 [(gogoproto.nullable) = false];
}

message MsgDelegateResponse {}
```

#### Undelegate

This message is used to unstake tokens from nft balance to a validator 

```protobuf

message MsgUndelegate {
  string creator = 1;
  string address = 2;
  string validator = 3;
  cosmos.base.v1beta1.Coin amount = 5 [(gogoproto.nullable) = false];
}

message MsgUndelegateResponse {}
```