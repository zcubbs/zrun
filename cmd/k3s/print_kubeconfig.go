package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/internal/configs"
	"os"
	"strings"
)

var (
	k3sServerUrl string // k3s server url
)

var printKubeconfigCmd = &cobra.Command{
	Use:     "print-kubeconfig",
	Aliases: []string{"kubeconfig", "kube", "config"},
	Short:   "print k3s kubeconfig",
	Long:    `print k3s kubeconfig`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		err := printKubeconfig(verbose)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func printKubeconfig(_ bool) error {
	kubeconfig := configs.Config.Kubeconfig.Path

	// read kubeconfig
	kubeconfigContent, err := readKubeconfig(kubeconfig)
	if err != nil {
		return err
	}

	if k3sServerUrl == "" {
		// get k3s url from env
		k3sUrl := os.Getenv("K3S_API_URL")
		if k3sUrl != "" {
			k3sServerUrl = k3sUrl
		} else {
			// get k3s url from config
			k3sServerUrl = "https://127.0.0.1:6443"
		}
	}

	// replace server url
	kubeconfigContent = replaceValueInString(kubeconfigContent,
		"https://127.0.0.1:6443", k3sServerUrl)

	// print kubeconfig
	fmt.Println(kubeconfigContent)

	return nil
}

func readKubeconfig(kubeconfig string) (string, error) {
	// read kubeconfig
	kubeconfigContent, err := os.ReadFile(kubeconfig)
	if err != nil {
		return "", err
	}
	return string(kubeconfigContent), nil
}

func replaceValueInString(str string, oldValue string, newValue string) string {
	return strings.ReplaceAll(str, oldValue, newValue)
}

func init() {
	printKubeconfigCmd.Flags().StringVarP(&k3sServerUrl, "url", "u", "", "k3s api server url. defaults to $K3S_API_URL or https://127.0.0.1:6443")

	Cmd.AddCommand(printKubeconfigCmd)
}
