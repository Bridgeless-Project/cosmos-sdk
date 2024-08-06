package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
)

func CmdNewAdmin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "newadmin  [from_key_or_address] [denom] [address] [rewards_per_period] [vesting_periods] [vesting_period]",
		Short: "Init a new admin with vesting params",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rewardPerPeriod, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			vestingPeriodsCount, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			vestingPeriod, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddAdmin(
				clientCtx.GetFromAddress().String(),
				args[1],
				args[2],
				rewardPerPeriod,
				vestingPeriodsCount,
				vestingPeriod,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
