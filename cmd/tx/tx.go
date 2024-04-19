package cmdtx

import (
	"os"

	"github.com/spf13/cobra"
	"piprim.net/gbcl/app"
	appaccount "piprim.net/gbcl/app/account"
	db "piprim.net/gbcl/app/db"
	tx "piprim.net/gbcl/app/tx"
	liberrors "piprim.net/gbcl/lib/errors"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

var txsCmd = &cobra.Command{
	Use:   "tx",
	Short: "Interact with transactions.",
}

func getTxCmdAdd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to database.",
		Run: func(cmd *cobra.Command, _ []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			tx := tx.New(appaccount.New(from), appaccount.New(to), value, data)

			state, err := db.NewStateFromDisk()
			if err != nil {
				liberrors.HandleError(err)
				os.Exit(1)
			}
			defer state.Close()

			err = state.AddTx(tx)
			if err != nil {
				liberrors.HandleError(err)
				os.Exit(1)
			}

			_, err = state.Persist()
			if err != nil {
				liberrors.HandleError(err)
				os.Exit(1)
			}

			app.Message("TX successfully added to the ledger.")
		},
	}

	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	liberrors.PanicOnErr(cmd.MarkFlagRequired(flagFrom))

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	liberrors.PanicOnErr(cmd.MarkFlagRequired(flagTo))

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	liberrors.PanicOnErr(cmd.MarkFlagRequired(flagValue))

	cmd.Flags().String(flagData, "", "Possible values: 'reward'")

	return cmd
}

func GetRootCmd() *cobra.Command {
	txsCmd.AddCommand(getTxCmdAdd())
	return txsCmd
}
