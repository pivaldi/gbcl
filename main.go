package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"piprim.net/gbcl/app"
	gcmd "piprim.net/gbcl/cmd"
	cmdbalance "piprim.net/gbcl/cmd/balance"
	cmdnode "piprim.net/gbcl/cmd/node"
	cmdtx "piprim.net/gbcl/cmd/tx"
	cmdversion "piprim.net/gbcl/cmd/version"
	liberrors "piprim.net/gbcl/lib/errors"
)

var gbclCmd *cobra.Command

func main() {
	initCmd()

	err := gbclCmd.Execute()
	liberrors.HandleErrorExit(err)
}

func initCmd() {
	gbclCmd = &cobra.Command{
		Use:     app.Name,
		Short:   app.ShortDescription,
		Version: app.GetVersion(),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			dataDir, _ := cmd.Flags().GetString(gcmd.FlagDataDir)
			err := app.Init(dataDir)

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	gcmd.AddPersistentFlags(gbclCmd)

	gbclCmd.AddCommand(cmdversion.GetRootCmd())
	gbclCmd.AddCommand(cmdbalance.GetRootCmd())
	gbclCmd.AddCommand(cmdnode.GetRootCmd())
	gbclCmd.AddCommand(cmdtx.GetRootCmd())
}
