// Package certmanager
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package certmanager

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/configs"
	helmPkg "github.com/zcubbs/zrun/helm"
	"helm.sh/helm/v3/pkg/cli/values"
	"k8s.io/utils/strings/slices"
	"log"
)

var (
	chartVersion string
	options      values.Options
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install cert-manager Chart",
	Long:  `install cert-manager Chart. Note: requires helm`,
	Run: func(cmd *cobra.Command, args []string) {
		err := installChart()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func installChart() error {
	fmt.Println("installing cert-manager Chart")
	kubeconfig := configs.Config.Kubeconfig.Path

	// chack if options.Values contains "installCRDs"
	// if not, add it
	if !slices.Contains(options.Values, "installCRDs=true") {
		options.Values = append(options.Values, "installCRDs=true")
	}

	verbose := Cmd.Flag("verbose").Value.String() == "true"

	helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		RepoName:     "cert-manager",
		RepoUrl:      "jetstack",
		ChartName:    "https://charts.jetstack.io",
		Namespace:    "cert-manager",
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
	})

	return nil
}

func init() {
	// parse flags
	install.Flags().StringVar(&chartVersion, "version", "", "chart version")
	install.Flags().StringArrayVar(&options.Values, "set", nil, "chart values")

	Cmd.AddCommand(install)
}
