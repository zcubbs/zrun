// Package k8s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k8s

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/internal/configs"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
)

var newNamespaces []string

// createNamespace represents the list command
var createNamespaceCmd = &cobra.Command{
	Use:   "create-ns",
	Short: "Show k8s create-ns list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("create k8s namespace")
		util.Must(
			util.RunTask(func() error {
				err := createNamespace(
					configs.Config.Kubeconfig.Path,
					newNamespaces,
				)
				return err
			}, true))
	},
}

func createNamespace(kubeconfig string, namespaces []string) error {
	err := kubectl.CreateNamespace(
		kubeconfig,
		namespaces,
	)
	if err != nil {
		return fmt.Errorf("couldn't create namespace\n %v", err)
	}

	return nil
}

func ExecuteCreateNamespaceCmd(kubeconfig, namespace string) error {
	err := kubectl.CreateNamespace(
		kubeconfig,
		[]string{namespace},
	)
	if err != nil {
		return fmt.Errorf("couldn't create namespace\n %v", err)
	}

	return nil
}

func init() {
	createNamespaceCmd.Flags().StringSliceVarP(&newNamespaces, "namespace", "n", nil, "namespace value")
	if err := createNamespaceCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(createNamespaceCmd)
}
