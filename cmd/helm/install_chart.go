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
	"log"
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

		fmt.Println("kubeconfig: ", kubeconfig)
		fmt.Println("repoName: ", repoName)
		fmt.Println("repoUrl: ", repoUrl)
		fmt.Println("chartName: ", chartName)
		fmt.Println("namespace: ", namespace)
		fmt.Println("chartVersion: ", chartVersion)
		fmt.Printf("chartValues: %+v", chartValues)

		// Execute Command
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		ExecuteInstallChartCmd(helm.InstallChartOptions{
			Kubeconfig:   kubeconfig,
			RepoName:     repoName,
			RepoUrl:      repoUrl,
			ChartName:    chartName,
			Namespace:    namespace,
			ChartVersion: chartVersion,
			ChartValues:  chartValues,
			Debug:        verbose,
		})
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
		log.Println(err)
	}
	if err := installChart.MarkFlagRequired("repo-url"); err != nil {
		log.Println(err)
	}
	if err := installChart.MarkFlagRequired("chart-name"); err != nil {
		log.Println(err)
	}
	if err := installChart.MarkFlagRequired("namespace"); err != nil {
		log.Println(err)
	}

	Cmd.AddCommand(installChart)
}

func ExecuteInstallChartCmd(options helm.InstallChartOptions) {
	// Add helm repo
	helm.RepoAdd(options.RepoName, options.RepoUrl)
	// Update charts from the helm repo
	helm.RepoUpdate()
	// Create Namespace
	k8s.ExecuteCreateNamespaceCmd(options.Kubeconfig, options.Namespace)
	// Install charts
	helm.InstallChart(options)
	// List helm releases
	ExecuteHelmListCmd()
}
