package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/accumulator/keeper"
	"github.com/cosmos/cosmos-sdk/x/accumulator/types"
	"github.com/spf13/cobra"
)

func CmdDerivePoolAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "derive-pool",
		Short: "Derive pool addressess",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Pool %s address: %s\n", types.AdminPoolName, keeper.GetPoolAddress(types.AdminPoolName).String())
			fmt.Printf("Pool %s address: %s\n", types.ValidatorPoolName, keeper.GetPoolAddress(types.ValidatorPoolName).String())
			fmt.Printf("Pool %s address: %s\n", types.NFTPoolName, keeper.GetPoolAddress(types.NFTPoolName).String())

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
