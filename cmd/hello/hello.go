// Package hello
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package hello

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "hello",
	Short: "hello is used for test purposes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		must.Succeed(
			progress.RunTask(func() error {
				style.PrintColoredHeader("Hello")
				style.PrintInfo("world!")
				return nil
			}, verbose))
	},
}
