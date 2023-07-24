// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "traefik",
	Short: "traefik Commands",
	Long:  `This command manages traefik.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}
