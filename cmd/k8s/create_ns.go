package k8s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/configs"
	"github.com/zcubbs/zrun/kubectl"
	"log"
)

var namespace string

// createNamespace represents the list command
var createNamespace = &cobra.Command{
	Use:   "create-ns",
	Short: "Show k8s create-ns list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteCreateNamespaceCmd(
			configs.Config.Kubeconfig.Path,
			namespace,
		)
	},
}

func ExecuteCreateNamespaceCmd(kubeconfig, namespace string) {
	err := kubectl.CreateNamespace(
		kubeconfig,
		namespace,
	)
	if err != nil {
		log.Fatalf("couldn't create namespace\n %v", err)
	}
}

func init() {
	createNamespace.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace value")
	if err := createNamespace.MarkFlagRequired("namespace"); err != nil {
		log.Println(err)
	}

	Cmd.AddCommand(createNamespace)
}
