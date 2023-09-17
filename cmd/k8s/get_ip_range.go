package k8s

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/zrun/internal/configs"
)

var getIpRangeCmd = &cobra.Command{
	Use:     "get-ip-range",
	Aliases: []string{"ip"},
	Short:   "get ip range",
	Long:    `get ip range`,
	Run: func(cmd *cobra.Command, args []string) {
		err := getIpRange(cmd.Context())
		if err != nil {
			fmt.Println(err)
		}
	},
}

func getIpRange(ctx context.Context) error {
	kubeconfig := configs.Config.Kubeconfig.Path
	ips, err := kubernetes.GetServiceCIDR(ctx, kubeconfig)
	if err != nil {
		return err
	}

	fmt.Println(ips)
	return nil
}

func init() {
	Cmd.AddCommand(getIpRangeCmd)
}
