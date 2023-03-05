// Package install
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package install

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install tools and apps",
	Long:  `This command installs tools and apps.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}
