// Package k9s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k9s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/k3s"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

// upgrade represents the list command
var k9s = &cobra.Command{
	Use:   "install",
	Short: "install k9s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("install k9s")

		must.Succeed(
			progress.RunTask(func() error {
				err := k3s.InstallK9s(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func init() {
	Cmd.AddCommand(k9s)
}
