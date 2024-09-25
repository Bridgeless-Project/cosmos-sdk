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

	var tallyParams v1.TallyParams
	paramSpace.Get(ctx, v1.ParamStoreKeyTallyParams, &tallyParams)

	tallyParams.DepositLockingPeriod = 7200 * time.Second

	paramSpace.Set(ctx, v1.ParamStoreKeyTallyParams, &tallyParams)
	return nil
}

func migrateVotingParams(ctx sdk.Context, paramSpace types.ParamSubspace) error {
	var oldVotingParams OldVotingParams
	paramSpace.Get(ctx, v1.ParamStoreKeyTallyParams, &oldVotingParams)

	vp := oldVotingParams.GetVotingPeriod()
	votingParams := v1.VotingParams{
		VotingPeriod: *vp,
	}

	paramSpace.Set(ctx, v1.ParamStoreKeyTallyParams, &votingParams)
	return nil
}

func MigrateStore(ctx sdk.Context, paramSpace types.ParamSubspace) error {
	if err := migrateTallyParams(ctx, paramSpace); err != nil {
		return err
	}

	if err := migrateVotingParams(ctx, paramSpace); err != nil {
		return err
	}

	return nil
}
