
# Client

## Query


### Admins
This query returns list of admins

#### CLI

```sh
simd noaccumulator admins --flags"
```

Example:
```
simd-cored query noaccumulator admins
```
Example Output:

```
admins: ["bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4"]
pagination:
  next_key: null
  total: "0"
```

#### HTTP

This query returns list of admins
```
cosmos/cosmos-sdk/accumulator/admins
```

Example:
```
http://localhost:1317/cosmos/cosmos-sdk/accumulator/admins
```
Example Output
```
{
  "admins": ["bridge103n4cmjt2je8nqcxg9y9desyhy6m57u52kkuc4"], 
  "pagination": {
    "next_key": null,
    "total": "0"
  }
}
```


## Message

### AddAdmin

This message is used to add a new admin to the `accumulator` module.

```protobuf
message MsgAddAdmin {
  string creator = 1;
  string address = 2;
  cosmos.base.v1beta1.Coin reward_per_period = 3 [(gogoproto.nullable) = false];
  int64 vesting_periods_count = 4;
  string denom = 5;
  int64 vesting_period = 6;
}

message MsgAddAdminResponse {}
```

