// Package k9s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k9s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/bash"
	"github.com/zcubbs/zrun/configs"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "k9s",
	Short: "Run K9s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := bash.ExecuteCmd(
			"k9s",
			"--kubeconfig",
			configs.Config.Kubeconfig.Path,
		)

		if err != nil {
			panic(err)
		}
	},
}
