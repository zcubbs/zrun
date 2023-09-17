// Package k8s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k8s

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
	"github.com/zcubbs/zrun/internal/configs"
)

var (
	crSecretName       string
	crSecretServer     string
	crSecretUsername   string
	crSecretPassword   string
	crSecretEmail      string
	crSecretNamespaces []string
	crSecretReplace    bool
)

// createContainerRegistryCmd represents the create k8s secret command
var createContainerRegistryCmd = &cobra.Command{
	Use:   "create-cr-secret",
	Short: "Create k8s container registry secret",
	Long:  `Create k8s container registry secret`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("add container registry secret")

		must.Succeed(
			progress.RunTask(func() error {
				err := createContainerRegistrySecret(cmd.Context(), verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func createContainerRegistrySecret(ctx context.Context, debug bool) error {
	kubeconfig := configs.Config.Kubeconfig.Path

	err := kubernetes.CreateContainerRegistrySecret(
		ctx,
		kubeconfig,
		kubernetes.ContainerRegistrySecret{
			Name:     crSecretName,
			Server:   crSecretServer,
			Username: crSecretUsername,
			Password: crSecretPassword,
			Email:    crSecretEmail,
		},
		crSecretNamespaces,
		crSecretReplace,
		debug,
	)

	if err != nil {
		return err
	}

	return nil
}

func init() {
	createContainerRegistryCmd.Flags().StringSliceVar(&crSecretNamespaces, "namespace", nil, "namespace value")
	createContainerRegistryCmd.Flags().StringVar(&crSecretName, "name", "", "name value")
	createContainerRegistryCmd.Flags().StringVar(&crSecretServer, "server", "", "server value")
	createContainerRegistryCmd.Flags().StringVar(&crSecretUsername, "username", "", "username value")
	createContainerRegistryCmd.Flags().StringVar(&crSecretPassword, "password", "", "password value")
	createContainerRegistryCmd.Flags().StringVar(&crSecretEmail, "email", "", "email value")
	createContainerRegistryCmd.Flags().BoolVar(&crSecretReplace, "replace", false, "replace value")

	_ = createContainerRegistryCmd.MarkFlagRequired("namespace")
	_ = createContainerRegistryCmd.MarkFlagRequired("name")
	_ = createContainerRegistryCmd.MarkFlagRequired("server")
	_ = createContainerRegistryCmd.MarkFlagRequired("username")
	_ = createContainerRegistryCmd.MarkFlagRequired("password")
	_ = createContainerRegistryCmd.MarkFlagRequired("email")

	Cmd.AddCommand(createContainerRegistryCmd)
}
