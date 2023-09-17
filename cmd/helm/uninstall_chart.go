// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/zrun/internal/configs"
	"os"
)

// uninstallChart represents the list command
var uninstallChart = &cobra.Command{
	Use:   "uninstall-chart",
	Short: "list all helm releases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		kubeconfig = configs.Config.Kubeconfig.Path
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		// Install charts
		err = helm.UninstallChart(kubeconfig, chartName, namespace, verbose)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	uninstallChart.Flags().StringVarP(&chartName, "chart-name", "c", "", "Helm chart name")
	uninstallChart.Flags().StringVarP(&namespace, "namespace", "n", "", "Helm chart namespace")

	if err := uninstallChart.MarkFlagRequired("chart-name"); err != nil {
		fmt.Println(err)
	}
	if err := uninstallChart.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(uninstallChart)
}
