// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/k8s"
	"github.com/zcubbs/zrun/configs"
	"github.com/zcubbs/zrun/helm"
	"helm.sh/helm/v3/pkg/cli/values"
	"os"
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
			fmt.Println(err)
			os.Exit(1)
		}
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
	fmt.Println("-------------------------------------------")
	fmt.Printf("installing '%s' helm Chart ...\n", options.ChartName)
	if options.Debug {
		fmt.Println("kubeconfig: ", options.Kubeconfig)
		fmt.Println("repoName: ", options.RepoName)
		fmt.Println("repoUrl: ", options.RepoUrl)
		fmt.Println("chartName: ", options.ChartName)
		fmt.Println("namespace: ", options.Namespace)
		fmt.Println("chartVersion: ", options.ChartVersion)
		fmt.Printf("chartValues: %+v", options.ChartValues)
		fmt.Printf("Helm options: %+v\n", options)
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
	k8s.ExecuteCreateNamespaceCmd(options.Kubeconfig, options.Namespace)

	// Install charts
	err = helm.InstallChart(options)
	if err != nil {
		return err
	}

	// List helm releases
	err = ExecuteHelmListCmd()
	if err != nil {
		return err
	}

	return nil
}
