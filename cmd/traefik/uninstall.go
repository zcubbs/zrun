package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/style"
	"github.com/zcubbs/zrun/internal/configs"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall traefik",
	Long:  `uninstall traefik`,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("uninstall traefik")
		err := uninstall()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func uninstall() error {
	kubeconfig := configs.Config.Kubeconfig.Path
	err := helm.UninstallChart(
		kubeconfig,
		"traefik",
		traefikNamespace,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to uninstall traefik: %w", err)
	}

	return nil
}

func init() {
	Cmd.AddCommand(uninstallCmd)
}
