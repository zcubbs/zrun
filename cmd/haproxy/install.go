// Package haproxy
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package haproxy

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/must"
	xos "github.com/zcubbs/x/os"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install haproxy package",
	Long:  `install haproxy package. Note: tested on Ubuntu 20.04+`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("install haproxy")

		must.Succeed(
			progress.RunTask(func() error {
				err := installHaproxy(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func installHaproxy(verbose bool) error {
	err := xos.Install("haproxy")
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Cmd.AddCommand(install)
}
