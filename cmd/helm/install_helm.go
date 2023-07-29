// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/helm"
	"os"
)

// installHelm represents the list command
var installHelm = &cobra.Command{
	Use:   "install-helm",
	Short: "install helm CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := ExecuteInstallHelmCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	Cmd.AddCommand(installHelm)
}

func ExecuteInstallHelmCmd() error {
	fmt.Println("-------------------------------------------")
	fmt.Println("installing helm...")
	// check if helm is installed
	installed, err := helm.IsHelmInstalled()
	if err != nil {
		return err
	}

	if installed {
		fmt.Println("helm is already installed")
		return nil
	}

	err = helm.Install(false)
	if err != nil {
		return err
	}
	return nil
}
