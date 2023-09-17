package argo

import (
	"fmt"
	"github.com/spf13/cobra"
	kubectl "github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/zrun/internal/configs"
)

const (
	EchoCmd = "kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath=\"{.data.password}\" | base64 -d; echo"
)

// printInitPasswd prints the argocd initial password
var printInitPasswd = &cobra.Command{
	Use:     "passwd",
	Aliases: []string{"password", "pass", "pwd"},
	Short:   "print argocd initial password",
	Long:    `print argocd initial password`,
	Run: func(cmd *cobra.Command, args []string) {
		err := PrintPasswd(true)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func PrintPasswd(_ bool) error {
	kubeconfig := configs.Config.Kubeconfig.Path
	secret, err := kubectl.GetSecret(kubeconfig, namespace, "argocd-initial-admin-secret")
	if err != nil {
		return err
	}

	fmt.Println(string(secret.Data["password"]))
	return nil
}

func init() {
	Cmd.AddCommand(printInitPasswd)
}
