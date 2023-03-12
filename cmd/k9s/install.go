// Package k9s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k9s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
)

// upgrade represents the list command
var k9s = &cobra.Command{
	Use:   "install",
	Short: "install k9s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := k3s.InstallK9s()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	Cmd.AddCommand(k9s)
}
