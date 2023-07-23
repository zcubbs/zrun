// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
	"log"
)

// uninstall represents the list command
var uninstall = &cobra.Command{
	Use:   "install",
	Short: "install k3s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := k3s.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(uninstall)
}
