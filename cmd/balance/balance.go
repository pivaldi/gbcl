package cmdbalance

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	appdb "piprim.net/gbcl/app/db"
	liberrors "piprim.net/gbcl/lib/errors"
)

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "Interact with balances.",
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(_ *cobra.Command, _ []string) {
		state, err := appdb.NewStateFromDisk()
		if err != nil {
			liberrors.HandleError(err)
			os.Exit(1)
		}
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
