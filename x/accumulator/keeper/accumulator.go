package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

// SetAdmin set a specific deposit in the store from its index
func (k BaseKeeper) SetAdmin(ctx sdk.Context, v types.Admin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AdminKeyPrefix))
	b := k.cdc.MustMarshal(&v)
	store.Set(types.AdminKey(v.Address), b)
}

// GetAdmin returns a Admin from its index
func (k BaseKeeper) GetAdmin(
	ctx sdk.Context,
	address string,
) (val types.Admin, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AdminKeyPrefix))
	b := store.Get(types.AdminKey(
		address,
	))

	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAdmin removes a Admin from the store
func (k BaseKeeper) RemoveAdmin(
	ctx sdk.Context,
	address string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AdminKeyPrefix))
	store.Delete(types.AdminKey(
		address,
	))
}

// GetAllAdmin returns all Admin
func (k BaseKeeper) GetAllAdmins(ctx sdk.Context) (list []types.Admin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AdminKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Admin
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllAdminsWithPagination returns all Admin with pagination
func (k BaseKeeper) GetAllAdminsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]types.Admin, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AdminKeyPrefix))

	var admins []types.Admin
	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
		var admin types.Admin

		k.cdc.MustUnmarshal(value, &admin)

		admins = append(admins, admin)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return admins, pageRes, nil
}
