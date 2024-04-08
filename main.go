package main

import (
	"os"

	"github.com/spf13/cobra"
	"piprim.net/gbcl/app"
	cmdbalance "piprim.net/gbcl/cmd/balance"
	cmdtx "piprim.net/gbcl/cmd/tx"
	cmdversion "piprim.net/gbcl/cmd/version"
	liberrors "piprim.net/gbcl/lib/errors"
)

var gbclCmd *cobra.Command

func main() {
	err := app.Init()
	if err != nil {
		handleError(err)
	}

	initCmd()

	err = gbclCmd.Execute()
	if err != nil {
		handleError(err)
	}
}

func initCmd() {
	config := app.GetConfig()

	gbclCmd = &cobra.Command{
		Use:     config.Name,
		Short:   config.ShortDescription,
		Version: app.GetVersion(),
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	gbclCmd.AddCommand(cmdversion.GetRootCmd())
	gbclCmd.AddCommand(cmdbalance.GetRootCmd())
	gbclCmd.AddCommand(cmdtx.GetRootCmd())
}

func handleError(err error) {
	liberrors.HandleError(err)
	os.Exit(1)
}
