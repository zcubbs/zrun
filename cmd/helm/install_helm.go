// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

// installHelm represents the list command
var installHelm = &cobra.Command{
	Use:   "install-helm",
	Short: "install helm CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("install helm cli")

		must.Succeed(
			progress.RunTask(func() error {
				err := ExecuteInstallHelmCmd()
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func init() {
	Cmd.AddCommand(installHelm)
}

func ExecuteInstallHelmCmd() error {
	verbose := Cmd.Flag("verbose").Value.String() == "true"
	// check if helm is installed
	installed, err := helm.IsHelmInstalled()
	if err != nil {
		return err
	}

	if installed && verbose {
		fmt.Println("helm is already installed")
		return nil
	}

	err = helm.Install(false)
	if err != nil {
		return err
	}
	return nil
}
