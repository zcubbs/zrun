// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/internal/configs"
	helmPkg "github.com/zcubbs/zrun/pkg/helm"
	"github.com/zcubbs/zrun/pkg/kubernetes"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
)

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
		style.PrintColoredHeader("install argocd")

		util.Must(
			util.RunTask(func() error {
				err := installChart(cmd.Context())
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func installChart(ctx context.Context) error {
	kubeconfig := configs.Config.Kubeconfig.Path
	verbose := Cmd.Flag("verbose").Value.String() == "true"

	options := helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		ChartName:    ArgocdString,
		RepoName:     ArgocdString,
		RepoUrl:      "https://argoproj.github.io/argo-helm",
		Namespace:    namespace,
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
	}

	err := helm.ExecuteInstallChartCmd(options)
	if err != nil {
		return err
	}

	err = kubernetes.IsDeploymentReady(
		ctx,
		kubeconfig,
		options.Namespace,
		[]string{
			ArgocdServerDeploymentName,
			ArgocdRepoServerDeploymentName,
			ArgocdRedisDeploymentName,
			ArgocdDexServerDeploymentName,
			ArgocdApplicationsetControllerDeploymentName,
			ArgocdNotificationsControllerDeploymentName,
		},
		options.Debug,
	)
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
