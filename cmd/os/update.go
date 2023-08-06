// Package os
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"github.com/spf13/cobra"
	zos "github.com/zcubbs/zrun/pkg/os"
)

// update represents the list command
var update = &cobra.Command{
	Use:   "update",
	Short: "OS Update",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := zos.Update()
		if err != nil {
			println(err.Error())
			panic("Could not Update OS")
		}
	},
}

func init() {
	Cmd.AddCommand(update)
}
