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
	zvault "github.com/zcubbs/zrun/vault"
	"os"
	"strings"
)

var (
	repositoryName string // repository name
	repositoryUrl  string // repository url

	repositoryUsername string // git username
	repositoryPassword string // git password

	repositoryType string // repository type

	repositoryUseEnv   bool // use env vars for credentials
	repositoryUseVault bool // use vault for credentials
)

const Git = "git"
const Helm = "helm"

// repository add a repository to ArgoCD
var repository = &cobra.Command{
	Use:   "repository",
	Short: "add repository to ArgoCD",
	Long:  `add repository to ArgoCD`,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader(fmt.Sprintf("add argocd %s repository", repositoryType))
		util.Must(
			util.RunTask(func() error {
				err := runAddRepo()
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func runAddRepo() error {
	// check if repository type not git or helm
	if repositoryType != Git && repositoryType != Helm {
		fmt.Println("error: repository type must be git or helm")
		return nil
	}

	urlValid := strings.HasPrefix(repositoryUrl, "http://") ||
		strings.HasPrefix(repositoryUrl, "https://")

	// check if url is valid
	if !urlValid {
		fmt.Printf("error: repository url must be valid url: %s", repositoryUrl)
		return nil
	}

	if repositoryType == Git {
		urlValid = strings.HasSuffix(repositoryUrl, ".git")
		if !urlValid {
			fmt.Printf("error: url must be valid git url: %s. %s",
				repositoryUrl,
				"example: https://example.com/example.git",
			)
			return nil
		}
	}

	// handle credentials
	err := handleCredentials()
	if err != nil {
		fmt.Printf("error: couldn't get credentials %s", err)
		return nil
	}

	// add repository
	err = addRepo()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func handleCredentials() error {
	// check if use vault and env vars
	if repositoryUseEnv && repositoryUseVault {
		return fmt.Errorf("error: can't use vault and env vars at the same time")
	}

	// check if use vault
	if repositoryUseVault {
		if repositoryUsername == "" {
			return fmt.Errorf("error: username vault key is empty")
		}

		if repositoryPassword == "" {
			return fmt.Errorf("error: password vault key is empty")
		}

		// get username and password from vault
		sv, err := zvault.NewSecretVault()
		if err != nil {
			return err
		}

		repositoryPassword, err = sv.GetSecret(repositoryPassword)
		if err != nil {
			return err
		}

		repositoryUsername, err = sv.GetSecret(repositoryUsername)
		if err != nil {
			return err
		}

		return nil
	}

	// check if use env vars
	if repositoryUseEnv {
		if repositoryUsername == "" {
			return fmt.Errorf("error: username env key is empty")
		}

		if repositoryPassword == "" {
			return fmt.Errorf("error: password env key is empty")
		}

		repositoryUsername = os.Getenv(repositoryUsername)
		repositoryPassword = os.Getenv(repositoryPassword)
		return nil
	}

	return nil
}

type argoRepo struct {
	Name      string
	Namespace string
	Url       string
	Username  string
	Password  string
	Type      string
}

var argoRepoTmpl = `---

apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  type: {{ .Type }}
{{- if eq .Type "helm" }}
  name: {{ .Name }}
{{- end }}
  url: {{ .Url }}
  username: {{ .Username }}
  password: {{ .Password }}
`

func addRepo() error {

	// create project
	repo := &argoRepo{
		Name:      repositoryName,
		Namespace: namespace,
		Url:       repositoryUrl,
		Username:  repositoryUsername,
		Password:  repositoryPassword,
		Type:      repositoryType,
	}

	// Apply template
	verbose := Cmd.Flag("verbose").Value.String() == "true"
	err := kubectl.ApplyManifest(argoRepoTmpl, repo, verbose)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// parse flags
	repository.Flags().StringVar(&repositoryName, "name", "", "repository name")
	repository.Flags().StringVar(&repositoryUrl, "url", "", "repository url")
	repository.Flags().StringVar(&repositoryUsername, "username", "", "repository username")
	repository.Flags().StringVar(&repositoryPassword, "password", "", "repository password")
	repository.Flags().StringVar(&repositoryType, "type", "repository type", "repository type")
	repository.Flags().BoolVar(&repositoryUseEnv, "use-env", false, "use env vars for credentials")
	repository.Flags().BoolVar(&repositoryUseVault, "use-vault", false, "use vault for credentials")

	// mark required flags
	err := repository.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}

	err = repository.MarkFlagRequired("url")
	if err != nil {
		fmt.Println(err)
	}

	err = repository.MarkFlagRequired("type")
	if err != nil {
		fmt.Println(err)
	}

	// add command to root
	Cmd.AddCommand(repository)
}
