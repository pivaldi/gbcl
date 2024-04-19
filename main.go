package main

import (
	"github.com/spf13/cobra"
	"piprim.net/gbcl/app"
	cmdbalance "piprim.net/gbcl/cmd/balance"
	cmdnode "piprim.net/gbcl/cmd/node"
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
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	gbclCmd.AddCommand(cmdversion.GetRootCmd())
	gbclCmd.AddCommand(cmdbalance.GetRootCmd())
	gbclCmd.AddCommand(cmdnode.GetRootCmd())
	// gbclCmd.AddCommand(cmdtx.GetRootCmd())
}
