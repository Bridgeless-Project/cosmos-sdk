package v046_26

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"time"
)

func migrateTallyParams(ctx sdk.Context, paramSpace types.ParamSubspace) error {
	ctx.Logger().Info(fmt.Sprintf("Performing v0.46.26 %s module migrations", types.ModuleName))

	var oldTallyParams v1.TallyParams
	paramSpace.Get(ctx, v1.ParamStoreKeyTallyParams, &oldTallyParams)

	oldTallyParams.DepositLockingPeriod = 7200 * time.Second
	paramSpace.Set(ctx, v1.ParamStoreKeyTallyParams, oldTallyParams)
	return nil
}

func MigrateStore(ctx sdk.Context, paramSpace types.ParamSubspace) error {
	return migrateTallyParams(ctx, paramSpace)
}
