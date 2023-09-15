package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config commands",
	Long:  `configure traefik`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	Cmd.AddCommand(configCmd)
}
