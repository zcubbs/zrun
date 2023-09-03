// Package haproxy
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package haproxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"os"
)

const (
	haproxyConfigFilePath = "/etc/haproxy/haproxy.cfg"
)

var (
	k3sDomain    string // k3s domain
	k3sApiDomain string // k3s api domain
	k3sNodeName  string // k3s node name
	k3sNodeIP    string // k3s node ip
)

// install represents the list command
var k3sSetupCmd = &cobra.Command{
	Use:   "k3s-setup",
	Short: "configure haproxy for k3s single node",
	Long:  `configure haproxy for k3s single node`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("install haproxy")

		util.Must(
			util.RunTask(func() error {
				err := configureHaproxyK3s(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

type K3sConfig struct {
	K3sDomain    string
	K3sApiDomain string
	K3sNodeName  string
	K3sNodeIP    string
}

func configureHaproxyK3s(verbose bool) error {
	k3sConfig := K3sConfig{
		K3sDomain:    k3sDomain,
		K3sApiDomain: k3sApiDomain,
		K3sNodeName:  k3sNodeName,
		K3sNodeIP:    k3sNodeIP,
	}

	configFileContent, err := kubectl.ApplyTmpl(haproxyK3sConfigTmpl, k3sConfig, verbose)
	if err != nil {
		return fmt.Errorf("failed to apply template \n %w", err)
	}

	// write tmp manifest
	err = os.WriteFile(haproxyConfigFilePath, configFileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write haproxy config file \n %w", err)
	}

	return nil
}

func init() {
	k3sSetupCmd.Flags().StringVarP(&k3sDomain, "k3s-domain", "d", "", "k3s domain")
	k3sSetupCmd.Flags().StringVarP(&k3sApiDomain, "k3s-api-domain", "a", "", "k3s api domain")
	k3sSetupCmd.Flags().StringVarP(&k3sNodeName, "k3s-node-name", "n", "", "k3s node name")
	k3sSetupCmd.Flags().StringVarP(&k3sNodeIP, "k3s-node-ip", "i", "", "k3s node ip")

	_ = k3sSetupCmd.MarkFlagRequired("k3s-domain")
	_ = k3sSetupCmd.MarkFlagRequired("k3s-api-domain")
	_ = k3sSetupCmd.MarkFlagRequired("k3s-node-name")
	_ = k3sSetupCmd.MarkFlagRequired("k3s-node-ip")

	Cmd.AddCommand(k3sSetupCmd)
}

var haproxyK3sConfigTmpl = `
defaults	
	log	global
	mode	http
	option	httplog
	option	dontlognull
	timeout http-request 20s
	timeout connect 5000
	timeout client  50000 # ddos protection
	timeout server  50000 # stick-table type ip size 100k expire 30s store conn_cur
	timeout http-keep-alive 10s

frontend k3s_http	
	bind *:80
	mode tcp
	
	acl k3s hdr_end(host) -i {{ .K3sDomain }}

	use_backend 	k3s_ingress_http if k3s
	default_backend k3s_ingress_http

frontend k3s_https
	bind *:443
	mode tcp

    acl k3s_api hdr_end(host) -i {{ .K3sApiDomain }}
	acl k3s 	hdr_end(host) -i {{ .K3sDomain }}

	use_backend 	k3s_api 			if k3s_api
	use_backend 	k3s_ingress 		if k3s
	default_backend k3s_ingress_https

backend k3s_api
	balance roundrobin
	server {{ .K3sNodeName }} {{ .K3sNodeIP }}:6443 check check-ssl verify none inter 10000

backend k3s_ingress_http	
	balance roundrobin
	server {{ .K3sNodeName }} {{ .K3sNodeIP }}:8080 check

backend k3s_ingress_https
	balance roundrobin
	option ssl-hello-chk
	server {{ .K3sNodeName }} {{ .K3sNodeIP }}:444 send-proxy check
`
