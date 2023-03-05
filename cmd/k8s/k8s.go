package k8s

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes Management Commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
