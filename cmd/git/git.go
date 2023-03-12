package git

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "git",
	Short: "Git Helper Commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
