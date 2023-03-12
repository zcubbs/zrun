// Package awx
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package awx

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/bash"
	"github.com/zcubbs/zrun/configs"
	"github.com/zcubbs/zrun/helm"
	"log"
)

// upgrade represents the list command
var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall awx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := nukeOperator()
		if err != nil {
			log.Fatal(err)
		}

		err = nukeInstance(instanceTmpl, secretTmpl)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(uninstall)
}

func nukeOperator() error {
	kubeconfig := configs.Config.Kubeconfig.Path
	// Install charts
	helm.UninstallChart(kubeconfig, "awx-operator", "default")

	return nil
}

func nukeInstance(instanceTmplStr, secretTmplStr string) error {
	fmt.Println("deploying awx instance")

	var tmplData = map[string]string{
		"Namespace":    "awx",
		"InstanceName": "awx",
		"Password":     "admin",
	}

	err := applyTmpl(instanceTmplStr, tmplData)
	if err != nil {
		return err
	}

	err = applyTmpl(secretTmplStr, tmplData)
	if err != nil {
		return err
	}

	err = bash.ExecuteCmd("kubectl",
		"delete",
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
