// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/internal/configs"
	helmPkg "github.com/zcubbs/zrun/pkg/helm"
	"github.com/zcubbs/zrun/pkg/kubectl"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
	"os"
)

const (
	traefikNamespace = "traefik"
)

var (
	defaultArgs = [...]string{
		"--entrypoints.websecure.http.tls",
		"--global.sendanonymoususage=false",
	}
	insecureArgs = [...]string{
		"--serversTransport.insecureSkipVerify",
	}

	insecureForwardedHeadersArgs = [...]string{
		"--entrypoints.websecure.forwardedHeaders.insecure",
		"--entrypoints.web.forwardedHeaders.insecure",
	}

	proxyArgs = [...]string{
		"--entrypoints.websecure.proxyProtocol",
		"--entrypoints.websecure.proxyProtocol.insecure",
	}
)

var (
	chartVersion               string
	options                    values.Options
	additionalArgs             []string
	useDefaults                bool
	withInsecure               bool
	withForwardedHeaders       bool
	forwardedHeadersTrustedIps string
	withProxyProtocol          bool
	proxyProtocolTrustedIps    string
	ingressProvider            string
	endpointWeb                string
	endpointWebsecure          string
	debugLogs                  bool
	accessLogs                 bool
)

var (
	dnsProviderString string
	dnsResolver       string
	dnsResolverEmail  string
	dnsTz             string

	ovhEndpoint         string
	ovhEndpointEnvKey   = "OVH_ENDPOINT"
	ovhEndpointVaultKey = "ovhEndpoint"

	ovhAppKey         string
	ovhAppKeyEnvKey   = "OVH_APP_KEY"
	ovhAppKeyVaultKey = "ovhAppKey"

	ovhAppSecret         string
	ovhAppSecretEnvKey   = "OVH_APP_SECRET"
	ovhAppSecretVaultKey = "ovhAppSecret"

	ovhConsumerKey         string
	ovhConsumerKeyEnvKey   = "OVH_CONSUMER_KEY"
	ovhConsumerKeyVaultKey = "ovhConsumerKey"

	azureClientID         string
	azureClientIDEnvKey   = "AZURE_CLIENT_ID"
	azureClientIDVaultKey = "azureClientID"

	azureClientSecret         string
	azureClientSecretEnvKey   = "AZURE_CLIENT_SECRET"
	azureClientSecretVaultKey = "azureClientSecret"

	useVault bool
)

// installCmd represents the list command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install traefik Chart",
	Long:  `install traefik Chart. Note: requires helm`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("install traefik")

		util.Must(
			util.RunTask(func() error {
				err := installChart(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func installChart(verbose bool) error {
	kubeconfig := configs.Config.Kubeconfig.Path

	var additionalArgs []string

	if withInsecure {
		additionalArgs = append(additionalArgs, insecureArgs[:]...)

		if withForwardedHeaders {
			additionalArgs = append(additionalArgs, insecureForwardedHeadersArgs[:]...)
		}
	}

	if withProxyProtocol {
		additionalArgs = append(additionalArgs, proxyArgs[:]...)
	}

	// check if useDefaults is true, if so, use default values
	if useDefaults {
		additionalArgs = append(additionalArgs, defaultArgs[:]...)
	}

	// can't set both ingressProvider and dnsProvider
	if ingressProvider != "" && dnsProviderString != "" {
		return fmt.Errorf("can't set both ingressProvider and dnsProvider")
	}

	// check if ingressProvider is set, if so, use it
	if ingressProvider != "" {
		additionalArgs = append(additionalArgs, fmt.Sprintf(
			"%s=%s",
			"--providers.kubernetesIngress.ingressClass",
			ingressProvider,
		))
	}

	if dnsProviderString != "" {
		err := configureDNSChallengeVars()
		if err != nil {
			return fmt.Errorf("failed to configure dns challenge vars: %w", err)
		}
		err = createDnsSecret()
		if err != nil {
			return fmt.Errorf("failed to create dns secret: %w", err)
		}
	}

	valuesPath := "values.yaml"

	tv := traefikValues{
		AdditionalArguments: []string{},
		//AdditionalArguments: additionalArgs,
		DnsProvider:                        dnsProviderString,
		DnsResolver:                        dnsResolver,
		DnsResolverEmail:                   dnsResolverEmail,
		EnableDashboard:                    true,
		EnableAccessLog:                    accessLogs,
		DebugLog:                           debugLogs,
		EndpointsWeb:                       endpointWeb,
		EndpointsWebsecure:                 endpointWebsecure,
		ServersTransportInsecureSkipVerify: withInsecure,
		ForwardedHeaders:                   withForwardedHeaders,
		ForwardedHeadersTrustedIPs:         forwardedHeadersTrustedIps,
		ProxyProtocol:                      withProxyProtocol,
		ProxyProtocolTrustedIPs:            proxyProtocolTrustedIps,
		DnsTZ:                              dnsTz,
	}
	// create traefik values.yaml from template
	configFileContent, err := kubectl.ApplyTmpl(traefikValuesTmpl, tv, verbose)
	if err != nil {
		return fmt.Errorf("failed to apply template \n %w", err)
	}

	// write tmp manifest
	err = os.WriteFile(valuesPath, configFileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write traefik values.yaml \n %w", err)
	}

	options.ValueFiles = []string{valuesPath}

	err = helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		RepoName:     "traefik",
		RepoUrl:      "https://helm.traefik.io/traefik",
		ChartName:    "traefik",
		Namespace:    traefikNamespace,
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
		Upgrade:      true,
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// parse flags
	installCmd.Flags().BoolVarP(&debugLogs, "debug", "d", false, "debug logs")
	installCmd.Flags().BoolVarP(&accessLogs, "accesslog", "a", false, "access logs")
	installCmd.Flags().StringVar(&chartVersion, "version", "", "chart version")
	installCmd.Flags().StringSliceVar(&options.Values, "set", nil, "chart values")
	installCmd.Flags().StringSliceVar(&additionalArgs, "set-arg", nil, "chart values additional arguments")
	installCmd.Flags().BoolVar(&useDefaults, "defaults", false, "use default values")
	installCmd.Flags().BoolVar(&withInsecure, "insecure", false, "use insecure connection")
	installCmd.Flags().BoolVar(&withForwardedHeaders, "forwardedHeaders", false, "use insecure forwarded headers")
	installCmd.Flags().StringVar(&forwardedHeadersTrustedIps, "forwardedHeadersTrustedIPs", "", "forwarded headers trusted ips")
	installCmd.Flags().BoolVar(&withProxyProtocol, "proxy", false, "use proxy protocol")
	installCmd.Flags().StringVar(&proxyProtocolTrustedIps, "proxyProtocolTrustedIPs", "", "proxy protocol trusted ips")
	installCmd.Flags().StringVar(&ingressProvider, "ingressProvider", "", "ingress provider")
	installCmd.Flags().StringVar(&endpointWeb, "endpointWeb", "80", "endpoint web")
	installCmd.Flags().StringVar(&endpointWebsecure, "endpointWebsecure", "443", "endpoint websecure")

	installCmd.PersistentFlags().BoolVar(&useVault, "use-vault", false, "use vault")
	installCmd.PersistentFlags().StringVar(&dnsProviderString, "dns-provider", "", "dns provider")
	installCmd.PersistentFlags().StringVar(&dnsResolver, "dns-resolver", "letsencrypt", "dns resolver")
	installCmd.PersistentFlags().StringVar(&dnsTz, "dns-tz", "Europe/Paris", "dns tz")
	installCmd.PersistentFlags().StringVar(&dnsResolverEmail, "dns-resolver-email", "", "dns resolver email")

	Cmd.AddCommand(installCmd)
}

type traefikValues struct {
	AdditionalArguments                []string
	DnsProvider                        string
	DnsResolver                        string
	DnsResolverEmail                   string
	EnableDashboard                    bool
	EnableAccessLog                    bool
	DebugLog                           bool
	EndpointsWeb                       string
	EndpointsWebsecure                 string
	ServersTransportInsecureSkipVerify bool
	ForwardedHeaders                   bool
	ForwardedHeadersTrustedIPs         string
	ProxyProtocol                      bool
	ProxyProtocolTrustedIPs            string
	DnsTZ                              string
}

var traefikValuesTmpl = `
globalArguments:
  - "--global.checknewversion=false"
  - "--global.sendanonymoususage=false"
global:
  sendAnonymousUsage: false
  checkNewVersion: false
  log:
  {{- if .DebugLog }}
    level: DEBUG
  {{- else }}
    level: INFO
  {{- end }}
  accessLogs:	
  {{- if .EnableAccessLog }}	
    enabled: true
  {{- else }}
    enabled: false
  {{- end }}
service:
  enabled: true
  type: LoadBalancer
rbac:
  enabled: true
additionalArguments:
  {{- range $i, $arg := .AdditionalArguments }}
  - "{{ printf "%s" . }}"
  {{- end }}
  {{- if .ServersTransportInsecureSkipVerify }}
  - "--serversTransport.insecureSkipVerify"
  {{- end }}
  {{- if .ForwardedHeaders }}
  - "--entrypoints.websecure.forwardedHeaders.trustedIPs={{ .ForwardedHeadersTrustedIPs }}"
  - "--entrypoints.web.forwardedHeaders.trustedIPs={{ .ForwardedHeadersTrustedIPs }}"
  {{- end }}
  {{- if .ProxyProtocol }}
  - "--entrypoints.websecure.proxyProtocol.trustedIPs={{ .ProxyProtocolTrustedIPs }}"
  {{- end }}
  {{- if .DnsProvider }}
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge=true"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge.provider={{ .DnsProvider }}"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge.delayBeforeCheck=10"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.email={{ .DnsResolverEmail }}"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge=true"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge.provider={{ .DnsProvider }}"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge.delayBeforeCheck=10"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.email={{ .DnsResolverEmail }}"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
  {{- end }}
ports:
  websecure:
    tls:
      enabled: true
      certResolver: {{ .DnsResolver }}

persistence:
  enabled: true
  accessMode: ReadWriteOnce
  size: 128Mi
  path: /data
  annotations: { }

ingressRoute:
  dashboard:
    enabled: true

logs:
  general:
  {{- if .DebugLog }}
    level: DEBUG
  {{- else }}
	level: INFO
  {{- end }}
  access:
    enabled: true
pilot:
  enabled: false

deployment:
  initContainers:
    - name: volume-permissions
      image: busybox:1.31.1
      command: ["sh", "-c", "touch /data/acme.json; chmod -Rv 0600 /data/acme.json; cat /data/acme.json"]
      volumeMounts:
        - name: data
          mountPath: /data

envFrom:
  - secretRef:
      name: traefik-dns-account-credentials

`
