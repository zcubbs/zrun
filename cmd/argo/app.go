// Package argo
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package argo

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

var (
	appName          string   // app name
	appNamespace     string   // app namespace
	isHelm           bool     // is helm
	helmValueFiles   []string // helm value files
	project          string   // project
	repoURL          string   // repo url
	targetRevision   string   // target revision
	path             string   // path
	recurse          bool     // recurse
	createNamespace  bool     // create namespace
	prune            bool     // prune
	selfHeal         bool     // self heal
	allowEmpty       bool     // allow empty
	isOci            bool     // is oci
	ociChartName     string   // oci chart name
	ociRepoURL       string   // oci repo url
	ociChartRevision string   // oci chart revision
)

// Cmd represents the install command
var appCmd = &cobra.Command{
	Use:   "add-application",
	Short: "Manage ArgoCD Applications",
	Long:  `This command manages ArgoCD Applications`,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("add argocd application")
		verbose := cmd.Flag("verbose").Value.String() == "true"

		_ = progress.RunTask(func() error {
			must.Succeed(addApp(verbose))
			return nil
		}, true)
	},
}

type App struct {
	AppName          string
	AppNamespace     string
	IsHelm           bool
	HelmValueFiles   []string
	Project          string
	RepoURL          string
	TargetRevision   string
	Path             string
	Recurse          bool
	CreateNamespace  bool
	Prune            bool
	SelfHeal         bool
	AllowEmpty       bool
	ArgoNamespace    string
	IsOci            bool
	OciChartName     string
	OciChartRevision string
	OciRepoURL       string
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
      passCredentials: true
      valueFiles:
      {{- range .HelmValueFiles }}
        - {{ . }}
      {{- end }}
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

var argoAppOciTmpl = `---

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .AppName }}
  namespace: {{ .ArgoNamespace }}
spec:
  project: {{ .Project }}
  sources:
    - repoURL: {{ .OciRepoURL }}
      targetRevision: {{ .OciChartRevision }}
      chart: {{ .OciChartName }}
      helm:
        passCredentials: true
        valueFiles:
        {{- range .HelmValueFiles }}
          - {{ . }}
        {{- end }}	
    - repoURL: {{ .RepoURL }}
      targetRevision: {{ .TargetRevision }}
      ref: values
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
	if !isHelm && isOci {
		return fmt.Errorf("oci flag can only be used with helm charts. flag --helm flag is not set")
	}

	if isOci && ociChartName == "" {
		return fmt.Errorf("oci chart name cannot be empty, when --oci flag is set")
	}

	if (!isOci && !isHelm) && path == "" {
		return fmt.Errorf("path cannot be empty, when --helm flag is not set")
	}

	// create app
	a := &App{
		AppName:          appName,
		AppNamespace:     appNamespace,
		IsHelm:           isHelm,
		HelmValueFiles:   helmValueFiles,
		Project:          project,
		RepoURL:          repoURL,
		TargetRevision:   targetRevision,
		Path:             path,
		Recurse:          recurse,
		CreateNamespace:  createNamespace,
		Prune:            prune,
		SelfHeal:         selfHeal,
		AllowEmpty:       allowEmpty,
		ArgoNamespace:    namespace,
		IsOci:            isOci,
		OciChartName:     ociChartName,
		OciRepoURL:       ociRepoURL,
		OciChartRevision: ociChartRevision,
	}

	if isOci {
		// Apply template
		err := kubernetes.ApplyManifest(argoAppOciTmpl, a, verbose)
		if err != nil {
			return err
		}
		return nil
	}

	err := kubernetes.ApplyManifest(argoAppTmpl, a, verbose)
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
	appCmd.Flags().BoolVar(&isOci, "oci", false, "is oci")
	appCmd.Flags().StringVar(&ociChartName, "oci-chart", "", "oci chart name")
	appCmd.Flags().StringVar(&ociRepoURL, "oci-repo", "", "oci repo url")
	appCmd.Flags().StringVar(&ociChartRevision, "oci-revision", "HEAD", "oci chart revision")

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

	Cmd.AddCommand(appCmd)
}
