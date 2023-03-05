// Package os
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "os",
	Short: "OS Helper Commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("os called")
	},
}
