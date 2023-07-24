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

	helm.ExecuteInstallChartCmd(
		kubeconfig,
		"cert-manager",
		"jetstack",
		"https://charts.jetstack.io",
		"cert-manager",
		chartVersion,
		options,
	)

	return nil
}

func init() {
	// parse flags
	install.Flags().StringVar(&chartVersion, "version", "", "chart version")
	install.Flags().StringArrayVar(&options.Values, "set", nil, "chart values")

	Cmd.AddCommand(install)
}
