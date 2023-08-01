// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/kubectl"
	"github.com/zcubbs/zrun/style"
	"github.com/zcubbs/zrun/util"
)

var (
	appName         string   // app name
	appNamespace    string   // app namespace
	isHelm          bool     // is helm
	helmValueFiles  []string // helm value files
	project         string   // project
	repoURL         string   // repo url
	targetRevision  string   // target revision
	path            string   // path
	recurse         bool     // recurse
	createNamespace bool     // create namespace
	prune           bool     // prune
	selfHeal        bool     // self heal
	allowEmpty      bool     // allow empty
)

// Cmd represents the install command
var appCmd = &cobra.Command{
	Use:   "add-application",
	Short: "Manage ArgoCD Applications",
	Long:  `This command manages ArgoCD Applications`,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("add argocd application")
		verbose := cmd.Flag("verbose").Value.String() == "true"

		_ = util.RunTask(func() error {
			util.Must(addApp(verbose))
			return nil
		}, true)
	},
}

type ArgoApp struct {
	AppName         string
	AppNamespace    string
	IsHelm          bool
	HelmValueFiles  []string
	Project         string
	RepoURL         string
	TargetRevision  string
	Path            string
	Recurse         bool
	CreateNamespace bool
	Prune           bool
	SelfHeal        bool
	AllowEmpty      bool
	ArgoNamespace   string
}

var argoAppTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .AppName }}
  namespace: {{ .ArgoNamespace }}
spec:
  project: {{ .Project }}
  source:
    repoURL: {{ .RepoURL }}
    targetRevision: {{ .TargetRevision }}
    path: {{ .Path }}
    {{ if .IsHelm }}
    helm:
      valueFiles:
      {{ range .HelmValueFiles }}
        - {{ . }}
      {{ end }}
    {{ else }}
    directory:
      recurse: {{ .Recurse }}
    {{ end }}
  destination:
    server: https://kubernetes.default.svc
    namespace: {{ .AppNamespace }}
  syncPolicy:
    syncOptions:
      - CreateNamespace={{ .CreateNamespace }}
    automated:
      prune: {{ .Prune }}
      selfHeal: {{ .SelfHeal }}
      allowEmpty: {{ .AllowEmpty }}
`

func addApp(verbose bool) error {
	// create app
	a := &ArgoApp{
		AppName:         appName,
		AppNamespace:    appNamespace,
		IsHelm:          isHelm,
		HelmValueFiles:  helmValueFiles,
		Project:         project,
		RepoURL:         repoURL,
		TargetRevision:  targetRevision,
		Path:            path,
		Recurse:         recurse,
		CreateNamespace: createNamespace,
		Prune:           prune,
		SelfHeal:        selfHeal,
		AllowEmpty:      allowEmpty,
		ArgoNamespace:   namespace,
	}

	err := kubectl.ApplyManifest(argoAppTmpl, a, verbose)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	appCmd.Flags().StringVar(&appName, "app-name", "", "app name")
	appCmd.Flags().StringVar(&appNamespace, "app-namespace", "", "app namespace")
	appCmd.Flags().BoolVar(&isHelm, "helm", false, "is helm")
	appCmd.Flags().StringSliceVar(&helmValueFiles, "helm-value-file", []string{"values.yaml"}, "helm value files")
	appCmd.Flags().StringVar(&project, "project", "default", "project")
	appCmd.Flags().StringVar(&repoURL, "repo", "", "repo url")
	appCmd.Flags().StringVar(&targetRevision, "target-revision", "HEAD", "target revision")
	appCmd.Flags().StringVar(&path, "path", "", "path")
	appCmd.Flags().BoolVar(&recurse, "recurse", true, "recurse")
	appCmd.Flags().BoolVar(&createNamespace, "create-namespace", true, "create namespace")
	appCmd.Flags().BoolVar(&prune, "auto-prune", true, "prune")
	appCmd.Flags().BoolVar(&selfHeal, "self-heal", true, "self heal")
	appCmd.Flags().BoolVar(&allowEmpty, "allow-empty", false, "allow empty")

	// mandatory flags
	err := appCmd.MarkFlagRequired("app-name")
	if err != nil {
		fmt.Println(err)
	}
	err = appCmd.MarkFlagRequired("app-namespace")
	if err != nil {
		fmt.Println(err)
	}
	err = appCmd.MarkFlagRequired("repo")
	if err != nil {
		fmt.Println(err)
	}
	err = appCmd.MarkFlagRequired("path")
	if err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(appCmd)
}
