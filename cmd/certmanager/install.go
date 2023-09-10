// Package certmanager
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package certmanager

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/internal/configs"
	helmPkg "github.com/zcubbs/zrun/pkg/helm"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
	"k8s.io/utils/strings/slices"
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
		style.PrintColoredHeader("install cert-manager")

		util.Must(
			util.RunTask(func() error {
				err := installChart()
				if err != nil {
					return err
				}
				return nil
			}, true))

	},
}

func installChart() error {
	kubeconfig := configs.Config.Kubeconfig.Path

	// chack if options.Values contains "installCRDs"
	// if not, add it
	if !slices.Contains(options.Values, "installCRDs=true") {
		options.Values = append(options.Values, "installCRDs=true")
	}

	verbose := Cmd.Flag("verbose").Value.String() == "true"

	err := helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		ChartName:    "cert-manager",
		RepoName:     "jetstack",
		RepoUrl:      "https://charts.jetstack.io",
		Namespace:    "cert-manager",
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
		Upgrade:      true,
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
