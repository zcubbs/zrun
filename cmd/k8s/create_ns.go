package k8s

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/internal/configs"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"os"
)

var namespace string

// createNamespace represents the list command
var createNamespace = &cobra.Command{
	Use:   "create-ns",
	Short: "Show k8s create-ns list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := ExecuteCreateNamespaceCmd(
			configs.Config.Kubeconfig.Path,
			namespace,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func ExecuteCreateNamespaceCmd(kubeconfig, namespace string) error {
	err := kubectl.CreateNamespace(
		kubeconfig,
		namespace,
	)
	if err != nil {
		return fmt.Errorf("couldn't create namespace\n %v", err)
	}

	return nil
}

func init() {
	createNamespace.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace value")
	if err := createNamespace.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(createNamespace)
}
