package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	"github.com/spf13/cobra"
)

func CmdQueryNFTsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner [address] --flags",
		Short: "Query the all nfts by owner address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryAllNFTsByOwner{
				Owner:      args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.GetAllNFTsByOwner(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllOwners() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owners",
		Short: "Query the all NFT holders",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryAllOwners{
				Pagination: pageReq,
			}

			res, err := queryClient.GetAllOwners(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
