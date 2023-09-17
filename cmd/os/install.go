// Package os
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"github.com/spf13/cobra"
	zos "github.com/zcubbs/x/os"
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "OS install packages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OS install packages", args)
		err := zos.Install(args...)
		if err != nil {
			println(err.Error())
			panic("Could not install packages")
		}
	},
}

func init() {
	Cmd.AddCommand(install)
}
