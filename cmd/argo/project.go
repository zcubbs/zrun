// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/kubectl"
	"log"
)

var (
	projectName string // project name
)

// install represents the list command
var project = &cobra.Command{
	Use:   "project",
	Short: "add project to ArgoCD",
	Long:  `add project to ArgoCD`,
	Run: func(cmd *cobra.Command, args []string) {
		err := addProject()
		if err != nil {
			log.Fatal(err)
		}
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
		Namespace: namespace,
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
	project.Flags().StringVar(&projectName, "name", "", "project name")

	// make flags required
	err := project.MarkFlagRequired("name")
	if err != nil {
		log.Println(err)
	}

	// add command to root
	Cmd.AddCommand(project)
}
