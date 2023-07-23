// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"bytes"
	"fmt"
	"github.com/zcubbs/zrun/bash"
	"os"
	"text/template"
)

const InstallScript = "/tmp/k3s-install.sh"
const UninstallScript = "/usr/local/bin/k3s-uninstall.sh"

const ConfigTemplate = "config.tmpl"
const ConfigFileLocation = "/etc/rancher/k3s/k3s.yaml"

type Config struct {
	Disable                 []string
	TlsSan                  []string
	DataDir                 string
	DefaultLocalStoragePath string
	WriteKubeconfigMode     string
}

func Install(config Config) error {
	// prepare config file
	err := createConfigFileFromTemplate(config)
	if err != nil {
		return err
	}

	// curl -sfL https://get.k3s.io -o k3s-install.sh
	err = bash.ExecuteCmd(
		"curl",
		"https://get.k3s.io",
		"-o",
		InstallScript,
	)
	if err != nil {
		return err
	}

	// sh ./k3s-install.sh -s - --write-kubeconfig-mode 644
	err = os.Chmod("/tmp/k3s-install.sh", 0700)
	if err != nil {
		return err
	}

	_, err = bash.ExecuteScript(
		InstallScript,
		InstallScript,
		"-s",
		"-",
		"--write-kubeconfig-mode=644",
	)
	if err != nil {
		return err
	}

	return nil
}

func createConfigFileFromTemplate(config Config) error {
	tmpl, err := template.New("tmpManifest").Parse(ConfigTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return err
	}

	fmt.Println(buf.String())

	err = os.WriteFile(ConfigFileLocation, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func Uninstall() error {
	_, err := bash.ExecuteScript(
		UninstallScript,
		UninstallScript,
	)
	if err != nil {
		return err
	}

	return nil
}
