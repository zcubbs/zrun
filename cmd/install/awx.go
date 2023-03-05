// Package install
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package install

import (
	"fmt"
	"github.com/spf13/cobra"
)

// upgrade represents the list command
var awx = &cobra.Command{
	Use:   "awx",
	Short: "install awx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(template)
	},
}

func init() {
	Cmd.AddCommand(awx)
}

var template = `
---

apiVersion: awx.ansible.com/v1beta1
kind: AWX
metadata:
  name: %s
  namespace: %s
spec:
  service_type: ClusterIP
  ingress_type: none
`
