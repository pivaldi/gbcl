package cmdnode

import (
	"github.com/spf13/cobra"
	liberrors "piprim.net/gbcl/lib/errors"
	"piprim.net/gbcl/node"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Interact with node.",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launches the TBB node and its HTTP API.",
	Run: func(_ *cobra.Command, _ []string) {
		err := node.Run()
		liberrors.HandleError(err)
	},
}

func GetRootCmd() *cobra.Command {
	nodeCmd.AddCommand(runCmd)

	return nodeCmd
}
