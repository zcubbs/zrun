// Package awx
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package awx

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/bash"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/configs"
	"html/template"
	"log"
	"os"
)

// upgrade represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install awx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := installOperator()
		if err != nil {
			log.Fatal(err)
		}

		err = deployInstance(instanceTmpl, secretTmpl)
		if err != nil {
			log.Fatal(err)
		}

		err = configureAwx()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(install)
}

var instanceTmpl = `
---

apiVersion: awx.ansible.com/v1beta1
kind: AWX
metadata:
  name: {{ .InstanceName }}
  namespace: {{ .Namespace }}
spec:
  service_type: ClusterIP
  ingress_type: none
`

var secretTmpl = `
---
apiVersion: v1
kind: Secret
metadata:
  name: {{ .InstanceName }}-admin-password
  namespace: {{ .Namespace }}
stringData:
  password: {{ .Password }}
`

func installOperator() error {
	fmt.Println("installing awx operator")
	kubeconfig := configs.Config.Kubeconfig.Path
	helm.ExecuteInstallChartCmd(
		kubeconfig,
		"awx-operator",
		"awx-operator",
		"https://ansible.github.io/awx-operator/",
		"default",
		"",
		map[string]interface{}{},
	)

	return nil
}

func deployInstance(instanceTmplStr, secretTmplStr string) error {
	fmt.Println("deploying awx instance")
	var tmplData = map[string]string{
		"Namespace":    "default",
		"InstanceName": "awx",
		"Password":     "admin",
	}

	err := applyTmpl(instanceTmplStr, tmplData)
	if err != nil {
		return err
	}

	err = applyTmpl(secretTmplStr, tmplData)
	if err != nil {
		return err
	}

	err = bash.ExecuteCmd("kubectl",
		"get",
		"secret",
		"awx-admin-password",
		"-o",
		"jsonpath=\"{.data.password}\" | base64 --decode ; echo",
	)
	if err != nil {
		fmt.Println("failed to get awx admin password")
		return err
	}
	return nil
}

func configureAwx() error {
	return nil
}

func applyTmpl(tmplStr string, tmplData map[string]string) error {
	tmpl, err := template.New("tmpManifest").Parse(tmplStr)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return err
	}

	fmt.Println(buf.String())

	err = os.WriteFile("/tmp/tmpManifest.yaml", buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = bash.ExecuteCmd("kubectl", "apply", "-f", "/tmp/tmpManifest.yaml")
	if err != nil {
		return err
	}
	return nil
}
