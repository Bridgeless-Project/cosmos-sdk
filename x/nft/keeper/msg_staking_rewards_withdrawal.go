package keeper

import (
	"context"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	destributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
)

func (m msgServer) StakingWithdrawalRewards(goCtx context.Context, msg *types.MsgStakingWithdrawalRewards) (*types.MsgStakingWithdrawalRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Nft)
	if err != nil {
		return nil, err
	}

	amount, err := m.hooks.ProcessStakingRewardsWithdrawal(ctx, delegatorAddress, valAddr)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, destributiontypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Nft),
		),
	)

	return new(types.MsgStakingWithdrawalRewardsResponse), nil
}
