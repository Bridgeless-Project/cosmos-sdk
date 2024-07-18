package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	"github.com/spf13/cobra"
)

func CmdDerivePoolAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derive-pool",
		Short: "Derive pool addressess",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	fmt.Printf("Pool %s address: %s\n", "nft 1", keeper.GetPoolAddress("1").String())
	fmt.Printf("Pool %s address: %s\n", "nft 2", keeper.GetPoolAddress("2").String())
	fmt.Printf("Pool %s address: %s\n", "nft 3", keeper.GetPoolAddress("3").String())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
