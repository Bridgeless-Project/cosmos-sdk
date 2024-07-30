
# Client

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

