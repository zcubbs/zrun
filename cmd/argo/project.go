// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
)

var (
	projectName string // project name
)

// install represents the list command
var addProjectCmd = &cobra.Command{
	Use:   "add-project",
	Short: "add project to ArgoCD",
	Long:  `add project to ArgoCD`,
	Run: func(cmd *cobra.Command, args []string) {

		style.PrintColoredHeader("add argocd project")

		util.Must(
			util.RunTask(func() error {
				err := addProject()
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

var argoProjectTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  description: {{ .Name }}
  sourceRepos:
    - '*'
  clusterResourceWhitelist:
    - group: '*'
      kind: '*'
  destinations:
    - namespace: '*'
      server: https://kubernetes.default.svc

`

type argoAppProject struct {
	Name      string
	Namespace string
}

func addProject() error {
	// create project
	project := &argoAppProject{
		Name:      projectName,
		Namespace: "argo-cd",
	}

	// Apply template
	// use flag from parent command
	verbose := Cmd.Flag("verbose").Value.String() == "true"
	err := kubectl.ApplyManifest(argoProjectTmpl, project, verbose)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// parse flags
	addProjectCmd.Flags().StringVar(&projectName, "name", "", "project name")

	// make flags required
	err := addProjectCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}

	// add command to root
	Cmd.AddCommand(addProjectCmd)
}
