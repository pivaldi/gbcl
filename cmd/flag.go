package cmd

import (
	"github.com/spf13/cobra"
)

const FlagDataDir = "datadir"

func AddPersistentFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().String(FlagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")

	// err := cmd.MarkFlagRequired(FlagDataDir)
	// liberrors.HandleErrorExit(err)

	return cmd
}
