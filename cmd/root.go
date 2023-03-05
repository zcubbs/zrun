// Package cmd
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/viper"
	"github.com/zcubbs/zrun/cmd/config"
	"github.com/zcubbs/zrun/cmd/info"
	zos "github.com/zcubbs/zrun/cmd/os"
	"github.com/zcubbs/zrun/defaults"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  "",
	}

	viperCommand = &cobra.Command{
		Run: func(c *cobra.Command, args []string) {
			fmt.Println(viper.GetString("Flag"))
		},
	}

	aboutCmd = &cobra.Command{
		Use:   "about",
		Short: "Print the info about crucible-cli",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			figure.NewColorFigure("ZRUN", "colossal", "red", true).Print()
			figure.NewColorFigure("zrun", "morse", "red", true).Print()
			fmt.Printf("Version: %s\n", defaults.Version)
			fmt.Println("<zrun> is a swiss army knife cli for devops engineers")
			fmt.Println("Copyright (c) 2023 zakaria.elbouwab (zcubbs)")
			fmt.Println("Repository: https://github.com/zcubbs/zrun")
		},
	}

	persistRootFlag bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addSubCommandPalettes() {
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(zos.Cmd)
	rootCmd.AddCommand(info.Cmd)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&persistRootFlag, "persist", "p", false, "Persist the CLI")
	rootCmd.AddCommand(aboutCmd)
	rootCmd.AddCommand(viperCommand)
	addSubCommandPalettes()
}
