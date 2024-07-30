# State

## NFT

Each NFT is represented in the blockchain by a unique address. This address is used to store a balance and to stake this balance to a validator. To generate an NFT address, the following code snippet can be used:

```go
sdk.Bech32ifyAddressBytes("bridge", address.Derive(authtypes.NewModuleAddress(acumulatortypes.ModuleName), []byte(strconv.FormatInt(id, 10))))
```

The NFT object contains specific information to describe a token:

```protobuf
message NFT {
  string address = 1;
  string owner = 2;
  string uri = 3;
  int64 vesting_period = 4;
  cosmos.base.v1beta1.Coin reward_per_period = 5 [(gogoproto.nullable) = false];
  int64 vesting_periods_count = 6;
  cosmos.base.v1beta1.Coin available_to_withdraw = 7 [(gogoproto.nullable) = false];
  int64 last_vesting_time = 8;
  int64 vesting_counter = 9;
  string denom = 10;
}
```

There are several commands to update the NFT store:

- To add a new NFT or update an existing NFT, use `SetNFT(ctx sdk.Context, v types.NFT)`.
- To remove an NFT, use `RemoveNFT(ctx sdk.Context, address string)`, where `address` is the unique address of the NFT.
- To get an NFT by address, use `GetNFT(ctx sdk.Context, address string)`, where `address` is the unique address of the NFT.
- To get all NFTs, there are two methods (with and without pagination):
    - `GetNFTsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]types.NFT, *query.PageResponse, error)`
    - `GetNFTs(ctx sdk.Context) (list []types.NFT)`

## Owner

The Owner struct matches the NFT holder (owner) address with the NFT address:

```protobuf
message Owner {
  string address = 1;
  string nft_address = 2;
}
```

There are several commands to update the owner store:

- To add a new NFT owner, use `SetOwnerNFT(ctx sdk.Context, owner, nftAddress string)`. This function creates a new branch in the store where the key is the owner address, and the leaf in this branch is the Owner object.
- To remove the NFT owner, use `RemoveOwnerNft(ctx sdk.Context, owner string, nftAddress string)`.
- To get all NFTs by owner address, use `GetAllNFTsByOwnerWithPagination(ctx sdk.Context, ownerAddress string, pagination *query.PageRequest) ([]types.NFT, *query.PageResponse, error)`.
- To get all addresses that hold any NFT `GetAllOwnersWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error)`