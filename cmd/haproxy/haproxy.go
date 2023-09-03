// Package haproxy
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package haproxy

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "haproxy",
	Short: "haproxy Commands",
	Long:  `This command manages haproxy.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
