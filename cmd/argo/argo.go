// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	ArgocdString                                 = "argocd"
	ArgocdServerDeploymentName                   = "argocd-server"
	ArgocdRepoServerDeploymentName               = "argocd-repo-server"
	ArgocdRedisDeploymentName                    = "argocd-redis"
	ArgocdDexServerDeploymentName                = "argocd-dex-server"
	ArgocdApplicationsetControllerDeploymentName = "argocd-applicationset-controller"
	ArgocdNotificationsControllerDeploymentName  = "argocd-notifications-controller"
)

var (
	namespace string // repository namespace
)

// Cmd represents the install command
var Cmd = &cobra.Command{
	Use:   "argo",
	Short: "ArgoCD Commands",
	Long:  `This command manages ArgoCD configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", ArgocdString, "namespace")
}
