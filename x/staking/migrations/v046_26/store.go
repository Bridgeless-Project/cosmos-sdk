package v046_26

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info(fmt.Sprintf("Performing v0.46.26 %s module migrations", types.ModuleName))

	store := ctx.KVStore(storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(cdc, iterator.Value())
		delegation.Timestamp = time.Time{}

		delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

		store.Set(
			types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()),
			types.MustMarshalDelegation(cdc, delegation),
		)
	}
	return nil
}
