package cmdversion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	app "piprim.net/gbcl/app"
	liberrors "piprim.net/gbcl/lib/errors"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version.",
	Run: func(_ *cobra.Command, _ []string) {
		err := app.Init("")
		if err != nil {
			liberrors.HandleError(err)
			os.Exit(1)
		}

		fmt.Println(app.GetVersion())
	},
}

func GetRootCmd() *cobra.Command {
	return versionCmd
}
