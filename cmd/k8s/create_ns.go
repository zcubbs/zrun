// Package k8s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k8s

import (
	"fmt"
	"github.com/spf13/cobra"
	kubectl "github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
	"github.com/zcubbs/zrun/internal/configs"
)

var newNamespaces []string

// createNamespace represents the list command
var createNamespaceCmd = &cobra.Command{
	Use:   "create-ns",
	Short: "Show k8s create-ns list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("create k8s namespace")
		must.Succeed(
			progress.RunTask(func() error {
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
