// Package k9s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k9s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/bash"
	"github.com/zcubbs/zrun/internal/configs"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "k9s",
	Short: "Run K9s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		err := bash.ExecuteCmd(
			"k9s",
			verbose,
			"--kubeconfig",
			configs.Config.Kubeconfig.Path,
		)

		if err != nil {
			panic(err)
		}
	},
}
