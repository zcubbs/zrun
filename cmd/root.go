// Package cmd
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zcubbs/zrun/cmd/argo"
	"github.com/zcubbs/zrun/cmd/awx"
	"github.com/zcubbs/zrun/cmd/certmanager"
	"github.com/zcubbs/zrun/cmd/config"
	"github.com/zcubbs/zrun/cmd/git"
	"github.com/zcubbs/zrun/cmd/hello"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/cmd/info"
	"github.com/zcubbs/zrun/cmd/k3s"
	"github.com/zcubbs/zrun/cmd/k8s"
	"github.com/zcubbs/zrun/cmd/k9s"
	zos "github.com/zcubbs/zrun/cmd/os"
	"github.com/zcubbs/zrun/cmd/traefik"
	"github.com/zcubbs/zrun/cmd/upgrade"
	"github.com/zcubbs/zrun/cmd/vault"
	"github.com/zcubbs/zrun/internal/defaults"
	"os"
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

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of zrun",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}

	aboutCmd = &cobra.Command{
		Use:   "about",
		Short: "Print the info about zrun",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			About()
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
	rootCmd.AddCommand(upgrade.Cmd)
	rootCmd.AddCommand(zos.Cmd)
	rootCmd.AddCommand(info.Cmd)
	rootCmd.AddCommand(hello.Cmd)
	rootCmd.AddCommand(awx.Cmd)
	rootCmd.AddCommand(helm.Cmd)
	rootCmd.AddCommand(k8s.Cmd)
	rootCmd.AddCommand(k3s.Cmd)
	rootCmd.AddCommand(k9s.Cmd)
	rootCmd.AddCommand(git.Cmd)
	rootCmd.AddCommand(vault.Cmd)
	rootCmd.AddCommand(certmanager.Cmd)
	rootCmd.AddCommand(traefik.Cmd)
	rootCmd.AddCommand(argo.Cmd)
}

func init() {
	// add verbose flag
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.AddCommand(aboutCmd)
	rootCmd.AddCommand(viperCommand)
	rootCmd.AddCommand(versionCmd)
	addSubCommandPalettes()
}

func About() {
	figure.NewColorFigure("ZRUN", "colossal", "red", true).Print()
	figure.NewColorFigure("zrun", "morse", "red", true).Print()
	fmt.Println(getFullVersion())
	fmt.Println(getDescription())
	fmt.Println("Copyright (c) 2023 zakaria.elbouwab (zcubbs)")
	fmt.Println("Repository: https://github.com/zcubbs/zrun")
}

func getVersion() string {
	return fmt.Sprintf("v%s", defaults.Version)
}

func getFullVersion() string {
	return fmt.Sprintf(`
Version: v%s
Commit: %s
Date: %s
`, defaults.Version, defaults.Commit, defaults.Date)
}

func getDescription() string {
	return fmt.Sprintf(`
/zrun/ is a comprehensive command-line interface (CLI) that provides
a range of functionalities from installing k3s,
managing Helm Deployments & Argocd applications to
Git operations and more...
`)
}
