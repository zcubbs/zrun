// Package awx
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package awx

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/internal/configs"
	"github.com/zcubbs/zrun/pkg/bash"
	helmPkg "github.com/zcubbs/zrun/pkg/helm"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
)

// upgrade represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install awx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		util.Must(installOperator(verbose))
		util.Must(deployInstance(instanceTmpl, secretTmpl, verbose))
		util.Must(configureAwx(verbose))
	},
}

func init() {
	Cmd.AddCommand(install)
}

var instanceTmpl = `
---

apiVersion: awx.ansible.com/v1beta1
kind: AWX
metadata:
  name: {{ .InstanceName }}
  namespace: {{ .Namespace }}
spec:
  service_type: ClusterIP
  ingress_type: none
`

var secretTmpl = `
---
apiVersion: v1
kind: Secret
metadata:
  name: {{ .InstanceName }}-admin-password
  namespace: {{ .Namespace }}
stringData:
  password: {{ .Password }}
`

func installOperator(verbose bool) error {
	kubeconfig := configs.Config.Kubeconfig.Path

	err := helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		RepoName:     "awx-operator",
		RepoUrl:      "https://ansible.github.io/awx-operator/",
		ChartName:    "awx-operator",
		Namespace:    "default",
		ChartVersion: "",
		ChartValues:  values.Options{},
		Debug:        verbose,
	})
	if err != nil {
		return err
	}

	return nil
}

func deployInstance(instanceTmplStr, secretTmplStr string, verbose bool) error {
	fmt.Println("deploying awx instance")
	var tmplData = map[string]string{
		"Namespace":    "default",
		"InstanceName": "awx",
		"Password":     "admin",
	}

	err := kubectl.ApplyManifest(instanceTmplStr, tmplData, false)
	if err != nil {
		return err
	}

	err = kubectl.ApplyManifest(secretTmplStr, tmplData, false)
	if err != nil {
		return err
	}

	err = bash.ExecuteCmd("kubectl",
		verbose,
		"get",
		"secret",
		"awx-admin-password",
		"-o",
		"jsonpath=\"{.data.password}\" | base64 --decode ; echo",
	)
	if err != nil {
		return fmt.Errorf("failed to get awx admin password: %w", err)
	}
	return nil
}

func configureAwx(_ bool) error {
	return nil
}
