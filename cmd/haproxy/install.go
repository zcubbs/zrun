// Package haproxy
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package haproxy

import (
	"github.com/spf13/cobra"
	xos "github.com/zcubbs/zrun/pkg/os"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install haproxy package",
	Long:  `install haproxy package. Note: tested on Ubuntu 20.04+`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("install haproxy")

		util.Must(
			util.RunTask(func() error {
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
