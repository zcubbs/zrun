// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/configs"
	helmPkg "github.com/zcubbs/zrun/helm"
	"helm.sh/helm/v3/pkg/cli/values"
	"os"
)

const ArgocdString = "argo-cd"

var (
	chartVersion string
	options      values.Options
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install argo-cd Chart",
	Long:  `install argo-cd Chart. Note: requires helm`,
	Run: func(cmd *cobra.Command, args []string) {
		err := installChart()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func installChart() error {
	kubeconfig := configs.Config.Kubeconfig.Path
	verbose := Cmd.Flag("verbose").Value.String() == "true"

	err := helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		ChartName:    ArgocdString,
		RepoName:     ArgocdString,
		RepoUrl:      "https://argoproj.github.io/argo-helm",
		Namespace:    ArgocdString,
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// parse flags
	install.Flags().StringVar(&chartVersion, "version", "", "chart version")
	install.Flags().StringArrayVar(&options.Values, "set", nil, "chart values")

	Cmd.AddCommand(install)
}
