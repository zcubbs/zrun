// Package awx
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package awx

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "awx",
	Short: "Awx Management Commands",
	Long:  `This command manages AWX Instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
