// Package info
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package info

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the info command
var Cmd = &cobra.Command{
	Use:   "info",
	Short: "Info is a palette that contains system info commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {

}
