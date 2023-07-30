// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
	"github.com/zcubbs/zrun/style"
	"github.com/zcubbs/zrun/util"
)

// uninstall represents the list command
var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall k3s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		style.PrintColoredHeader("uninstall k3s")
		util.Must(
			util.RunTask(func() error {
				err := k3s.Uninstall(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func init() {
	Cmd.AddCommand(uninstall)
}
