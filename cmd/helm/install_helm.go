package helm

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/helm"
	"log"
)

// installHelm represents the list command
var installHelm = &cobra.Command{
	Use:   "install-helm",
	Short: "install helm CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := ExecuteInstallHelmCmd()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(installHelm)
}

func ExecuteInstallHelmCmd() error {
	// Add helm repo
	err := helm.Install()
	if err != nil {
		return err
	}
	return nil
}
