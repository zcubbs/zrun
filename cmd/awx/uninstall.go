// Package awx
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package awx

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/bash"
	"github.com/zcubbs/x/helm"
	kubectl "github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/zrun/internal/configs"
)

// upgrade represents the list command
var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall awx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		must.Succeed(nukeOperator())
		must.Succeed(nukeInstance(instanceTmpl, secretTmpl))
	},
}

func init() {
	Cmd.AddCommand(uninstall)
}

func nukeOperator() error {
	kubeconfig := configs.Config.Kubeconfig.Path
	verbose := Cmd.Flag("verbose").Value.String() == "true"

	// Install charts
	err := helm.UninstallChart(kubeconfig, "awx-operator", "default", verbose)
	if err != nil {
		return err
	}

	return nil
}

func nukeInstance(instanceTmplStr, secretTmplStr string) error {
	fmt.Println("deploying awx instance")

	var tmplData = map[string]string{
		"Namespace":    "awx",
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

	verbose := Cmd.Flag("verbose").Value.String() == "true"

	err = bash.ExecuteCmd("kubectl",
		verbose,
		"secret",
		"awx-admin-password",
		"-o",
		"jsonpath=\"{.data.password}\" | base64 --decode ; echo",
	)
	if err != nil {
		fmt.Println("failed to get awx admin password")
		return err
	}
	return nil
}
