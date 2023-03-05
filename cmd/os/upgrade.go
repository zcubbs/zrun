package os

import (
	"github.com/spf13/cobra"
	zos "github.com/zcubbs/zrun/os"
)

// upgrade represents the list command
var upgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "OS Upgrade",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := zos.Upgrade()
		if err != nil {
			println(err.Error())
			panic("Could not Upgrade OS")
		}
	},
}

func init() {
	Cmd.AddCommand(upgrade)
}
