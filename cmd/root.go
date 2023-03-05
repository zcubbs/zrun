// Package cmd
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zcubbs/zrun/cmd/config"
	"github.com/zcubbs/zrun/cmd/info"
	zos "github.com/zcubbs/zrun/cmd/os"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crucible",
		Short: "This is a CLI for the Crucible Bot project",
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
			fmt.Printf(configs.Splash + "\n\n")
			fmt.Println("zrun")
			fmt.Printf("Version: %s\n", configs.Version)
			fmt.Println("This is a swiss army knife cli for devops engineers")
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
