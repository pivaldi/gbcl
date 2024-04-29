package cmdbalance

import (
	"fmt"

	"github.com/spf13/cobra"
	"piprim.net/gbcl/app"
	db "piprim.net/gbcl/app/db"
	"piprim.net/gbcl/cmd"
	liberrors "piprim.net/gbcl/lib/errors"
)

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "Interact with balances.",
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(c *cobra.Command, _ []string) {
		dataDir, err := c.Flags().GetString(cmd.FlagDataDir)
		liberrors.HandleErrorExit(err)

		err = app.Init(dataDir)
		liberrors.HandleErrorExit(err)

		state, err := db.NewStateFromDisk()
		liberrors.HandleErrorExit(err)
		defer state.Close()

		fmt.Println("Accounts balances")
		fmt.Println("_________________")
		fmt.Println("")

		if state != nil {
			for account, balance := range state.Balances {
				fmt.Printf("%s : %d\n", account, balance)
			}
		}
	},
}

func GetRootCmd() *cobra.Command {
	balancesCmd.AddCommand(balancesListCmd)
	return balancesCmd
}
