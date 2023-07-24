// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"github.com/zcubbs/zrun/bash"
	log "github.com/zcubbs/zrun/log"
)

func Install() error {
	// curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
	log.Infow("Installing helm",
		map[string]any{
			"command": "curl",
			"args":    "-fsSL -o  /tmp/get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3",
		},
	)
	err := bash.ExecuteCmd(
		"curl",
		"-fsSL",
		"-o",
		"/tmp/get_helm.sh",
		"https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3",
	)
	if err != nil {
		return err
	}

	log.Infow("Exec",
		map[string]any{
			"command": "ls",
			"args":    []string{"-l"},
		},
	)

	// chmod 700 get_helm.sh
	err = bash.Chmod("/tmp/get_helm.sh", 0700)
	if err != nil {
		return err
	}

	log.Infow("Exec script",
		map[string]any{
			"command": "/tmp/get_helm.sh",
		})
	// sh ./get_helm.sh
	_, err = bash.ExecuteScript(
		"/tmp/get_helm.sh",
		"/tmp/get_helm.sh",
	)
	if err != nil {
		return err
	}

	return nil
}
