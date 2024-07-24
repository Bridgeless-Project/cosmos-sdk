package accumulator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// this line is used by starport scaffolding # handler/msgServer
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddAdmin:
			res, err := msgServer.AddAdmin(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		}
		return nil, nil

		// this line is used by starport scaffolding # 1
	}
}
