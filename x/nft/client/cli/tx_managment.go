package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func CmdWithdrawal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawal [from_key_or_address] [nft_address]",
		Short: "Withdrawal allowed nft amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawal(
				clientCtx.GetFromAddress().String(),
				args[1],
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

func CmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [receiver] [nft_address]",
		Short: "send nft to another account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSend(
				clientCtx.GetFromAddress().String(),
				args[1],
				args[2],
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

func CmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [owner-address] [start-vesting-block]",
		Short: "mint nft",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			startVestingBlock, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			creatorAddress := clientCtx.GetFromAddress().String()
			if creatorAddress == "" {
				return fmt.Errorf("must provide creator address")
			}

			fmt.Println("creator address:", creatorAddress)

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := buildMintMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			msg.Creator = creatorAddress
			msg.Owner = args[0]
			msg.StartVestingBlock = startVestingBlock

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FlagSetMetadata())

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagMetadataUri)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func buildMintMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgMint, error) {
	metadataURI, err := fs.GetString(FlagMetadataUri)
	if err != nil {
		return txf, nil, errors.Wrap(err, "failed to get the metadata URI")
	}

	msg := &types.MsgMint{
		NftMetadataUri: metadataURI,
	}

	return txf, msg, nil
}
