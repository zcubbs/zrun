// Package vault
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package vault

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "vault",
	Short: "Vault Helper Commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
