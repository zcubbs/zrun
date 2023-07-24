// Package cmd
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/viper"
	"github.com/zcubbs/zrun/cmd/awx"
	"github.com/zcubbs/zrun/cmd/certmanager"
	"github.com/zcubbs/zrun/cmd/config"
	"github.com/zcubbs/zrun/cmd/git"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/cmd/info"
	"github.com/zcubbs/zrun/cmd/k3s"
	"github.com/zcubbs/zrun/cmd/k8s"
	"github.com/zcubbs/zrun/cmd/k9s"
	zos "github.com/zcubbs/zrun/cmd/os"
	"github.com/zcubbs/zrun/cmd/traefik"
	"github.com/zcubbs/zrun/cmd/vault"
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
		Short: "Print the info about zrun",
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
	rootCmd.AddCommand(awx.Cmd)
	rootCmd.AddCommand(helm.Cmd)
	rootCmd.AddCommand(k8s.Cmd)
	rootCmd.AddCommand(k3s.Cmd)
	rootCmd.AddCommand(k9s.Cmd)
	rootCmd.AddCommand(git.Cmd)
	rootCmd.AddCommand(vault.Cmd)
	rootCmd.AddCommand(certmanager.Cmd)
	rootCmd.AddCommand(traefik.Cmd)
}

func init() {
	rootCmd.AddCommand(aboutCmd)
	rootCmd.AddCommand(viperCommand)
	addSubCommandPalettes()
}
