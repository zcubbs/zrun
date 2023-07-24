// Package certmanager
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package certmanager

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "certmanager",
	Short: "cert-manager Commands",
	Long:  `This command manages cert-manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
