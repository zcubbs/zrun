// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
	"github.com/zcubbs/zrun/util"
)

// uninstall represents the list command
var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall k3s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		util.Must(k3s.Uninstall(verbose))

		fmt.Println("k3s uninstalled")
	},
}

func init() {
	Cmd.AddCommand(uninstall)
}
