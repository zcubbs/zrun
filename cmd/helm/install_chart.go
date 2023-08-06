// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/k8s"
	"github.com/zcubbs/zrun/internal/configs"
	"github.com/zcubbs/zrun/pkg/helm"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
)

var (
	kubeconfig   string
	repoName     string
	repoUrl      string
	chartName    string
	namespace    string
	chartVersion string
	chartValues  values.Options
)

// installChart represents the list command
var installChart = &cobra.Command{
	Use:   "install-chart",
	Short: "list all helm releases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig = configs.Config.Kubeconfig.Path
		// Execute Command
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader(fmt.Sprintf(
			"install '%s' helm Chart\n",
			chartName))

		util.Must(
			util.RunTask(func() error {
				err := ExecuteInstallChartCmd(helm.InstallChartOptions{
					Kubeconfig:   kubeconfig,
					RepoName:     repoName,
					RepoUrl:      repoUrl,
					ChartName:    chartName,
					Namespace:    namespace,
					ChartVersion: chartVersion,
					ChartValues:  chartValues,
					Debug:        verbose,
				})

				if err != nil {
					return err
				}

				return nil
			}, true))
	},
}

func init() {
	installChart.Flags().StringVar(&repoName, "repo-name", "", "helm repo name")
	installChart.Flags().StringVar(&repoUrl, "repo-url", "", "helm repo url")
	installChart.Flags().StringVar(&chartName, "chart-name", "", "chart name")
	installChart.Flags().StringVar(&namespace, "namespace", "", "chart namespace")
	installChart.Flags().StringVar(&chartVersion, "chart-version", "", "chart version")
	installChart.Flags().StringArrayVar(&chartValues.Values, "set", nil, "chart values")

	if err := installChart.MarkFlagRequired("repo-name"); err != nil {
		fmt.Println(err)
	}
	if err := installChart.MarkFlagRequired("repo-url"); err != nil {
		fmt.Println(err)
	}
	if err := installChart.MarkFlagRequired("chart-name"); err != nil {
		fmt.Println(err)
	}
	if err := installChart.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(installChart)
}

func ExecuteInstallChartCmd(options helm.InstallChartOptions) error {
	if options.Debug {
		style.PrintColoredDebug(fmt.Sprintf("kubeconfig: %s", options.Kubeconfig))
		style.PrintColoredDebug(fmt.Sprintf("repoName: %s", options.RepoName))
		style.PrintColoredDebug(fmt.Sprintf("repoUrl: %s", options.RepoUrl))
		style.PrintColoredDebug(fmt.Sprintf("chartName: %s", options.ChartName))
		style.PrintColoredDebug(fmt.Sprintf("namespace: %s", options.Namespace))
		style.PrintColoredDebug(fmt.Sprintf("chartVersion: %s", options.ChartVersion))
		style.PrintColoredDebug(fmt.Sprintf("chartValues: %+v", options.ChartValues))
		style.PrintColoredDebug(fmt.Sprintf("Helm options: %+v\n", options))
	}

	// Add helm repo
	err := helm.RepoAdd(options.RepoName, options.RepoUrl, options.Debug)
	if err != nil {
		return err
	}

	// Update charts from the helm repo
	err = helm.RepoUpdate(options.Debug)
	if err != nil {
		return err
	}

	// Create Namespace
	err = k8s.ExecuteCreateNamespaceCmd(options.Kubeconfig, options.Namespace)
	if err != nil {
		return err
	}

	// Install charts
	err = helm.InstallChart(options)
	if err != nil {
		return err
	}

	// List helm releases
	if options.Debug {
		err = ExecuteHelmListCmd()
		if err != nil {
			return err
		}
	}

	return nil
}
