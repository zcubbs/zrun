package helm

import (
	"github.com/zcubbs/zrun/bash"
)

func Install() error {
	// curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
	err := bash.ExecuteCmd(
		"curl",
		"https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3",
		"-fsSL",
		"-o",
		"get_helm.sh",
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

	// chmod 700 get_helm.sh
	err = bash.Chmod("get_helm.sh", 0700)
	if err != nil {
		return err
	}

	// sh ./get_helm.sh
	_, err = bash.ExecuteScript(
		"./get_helm.sh",
		"./get_helm.sh",
	)
	if err != nil {
		return err
	}

	return nil
}
