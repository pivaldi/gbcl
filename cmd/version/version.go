package cmdversion

import (
	"fmt"

	"github.com/spf13/cobra"
	app "piprim.net/gbcl/app"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version.",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(app.GetVersion())
	},
}

func GetRootCmd() *cobra.Command {
	return versionCmd
}
