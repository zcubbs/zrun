// Package upgrade
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/bash"
	"github.com/zcubbs/zrun/style"
	"github.com/zcubbs/zrun/util"
	"os"
)

const (
	InstallScriptURL = "https://raw.githubusercontent.com/zcubbs/zrun/main/scripts/install/install.sh"
	InstallScript    = "/tmp/zrun-install.sh"
)

// Cmd represents the os command
var Cmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade is used to upgrade zrun to latest version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		style.PrintColoredHeader("update zrun to latest version")
		util.Must(
			util.RunTask(func() error {
				err := upgrade(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func upgrade(verbose bool) error {
	fmt.Println("upgrade...")

	// curl -sfL https://raw.githubusercontent.com/zcubbs/zrun/main/scripts/install/install.sh -o zrun-install.sh
	err := bash.ExecuteCmd(
		"curl",
		verbose,
		InstallScriptURL,
		"-o",
		InstallScript,
	)
	if err != nil {
		return fmt.Errorf("error while executing curl %s -o %s\n%v\n",
			InstallScriptURL,
			InstallScript,
			err,
		)
	}

	// sh ./k3s-install.sh -s - --write-kubeconfig-mode 644
	err = os.Chmod(InstallScript, 0700)
	if err != nil {
		return fmt.Errorf("error while running %s \n%v",
			"chmod 0700 ./tmp/k3s-install.sh -s - --write-kubeconfig-mode 644",
			err)
	}

	ok, err := bash.ExecuteScript(
		InstallScript,
		verbose,
		InstallScript,
	)
	if !ok && err != nil {
		return fmt.Errorf("error while running %s \n%v",
			InstallScript,
			err)
	}

	return nil
}
