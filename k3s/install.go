// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"github.com/zcubbs/zrun/bash"
	"os"
)

func Install() error {
	// curl -sfL https://get.k3s.io -o k3s-install.sh
	err := bash.ExecuteCmd(
		"curl",
		"https://get.k3s.io",
		"-o",
		"k3s-install.sh",
	)
	if err != nil {
		return err
	}

	// ls -l
	err = bash.ExecuteCmd(
		"ls",
		"-l",
	)
	if err != nil {
		return err
	}

	// sh ./k3s-install.sh -s - --write-kubeconfig-mode 644
	err = os.Chmod("k3s-install.sh", 0700)
	if err != nil {
		return err
	}

	_, err = bash.ExecuteScript(
		"./k3s-install.sh",
		"./k3s-install.sh",
		"-s",
		"-",
		"--write-kubeconfig-mode=644",
	)
	if err != nil {
		return err
	}

	return nil
}
