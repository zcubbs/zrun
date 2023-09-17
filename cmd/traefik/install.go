// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
	helmPkg "github.com/zcubbs/x/helm"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
	"github.com/zcubbs/x/templates"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/internal/configs"
	"helm.sh/helm/v3/pkg/cli/values"
	"os"
)

const (
	traefikNamespace = "traefik"
)

var (
	chartVersion                 string
	options                      values.Options
	additionalArgs               []string
	useDefaults                  bool
	withInsecure                 bool
	withForwardedHeaders         bool
	withForwardedHeadersInsecure bool
	forwardedHeadersTrustedIps   string
	withProxyProtocol            bool
	withProxyProtocolInsecure    bool
	proxyProtocolTrustedIps      string
	ingressProvider              string
	endpointWeb                  string
	endpointWebsecure            string
	debugLogs                    bool
	accessLogs                   bool
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

	ovhAppSecret string
	/* #nosec */
	ovhAppSecretEnvKey   = "OVH_APP_SECRET"
	ovhAppSecretVaultKey = "ovhAppSecret"

	ovhConsumerKey         string
	ovhConsumerKeyEnvKey   = "OVH_CONSUMER_KEY"
	ovhConsumerKeyVaultKey = "ovhConsumerKey"

	azureClientID         string
	azureClientIDEnvKey   = "AZURE_CLIENT_ID"
	azureClientIDVaultKey = "azureClientID"

	azureClientSecret string
	/* #nosec */
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

		must.Succeed(
			progress.RunTask(func() error {
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

	// can't set both ingressProvider and dnsProvider
	if ingressProvider != "" && dnsProviderString != "" {
		return fmt.Errorf("can't set both ingressProvider and dnsProvider")
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
		AdditionalArguments:                []string{},
		IngressProvider:                    ingressProvider,
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
		ForwardedHeadersInsecure:           withForwardedHeadersInsecure,
		ForwardedHeadersTrustedIPs:         forwardedHeadersTrustedIps,
		ProxyProtocol:                      withProxyProtocol,
		ProxyProtocolInsecure:              withProxyProtocolInsecure,
		ProxyProtocolTrustedIPs:            proxyProtocolTrustedIps,
		DnsTZ:                              dnsTz,
	}
	// create traefik values.yaml from template
	configFileContent, err := templates.ApplyTmpl(traefikValuesTmpl, tv, verbose)
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
	installCmd.Flags().BoolVar(&withForwardedHeadersInsecure, "forwardedHeadersInsecure", false, "use insecure forwarded headers")
	installCmd.Flags().StringVar(&forwardedHeadersTrustedIps, "forwardedHeadersTrustedIPs", "", "forwarded headers trusted ips")
	installCmd.Flags().BoolVar(&withProxyProtocol, "proxy", false, "use proxy protocol")
	installCmd.Flags().BoolVar(&withProxyProtocolInsecure, "proxyInsecure", false, "use insecure proxy protocol")
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
	IngressProvider                    string
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
	ForwardedHeadersInsecure           bool
	ForwardedHeadersTrustedIPs         string
	ProxyProtocol                      bool
	ProxyProtocolInsecure              bool
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
  {{- if .ForwardedHeadersInsecure }}
  - "--entrypoints.websecure.forwardedHeaders.insecure"
  {{- end }}
  {{- if .ForwardedHeadersTrustedIPs }}
  - "--entrypoints.websecure.forwardedHeaders.trustedIPs=127.0.0.1/32,{{ .ForwardedHeadersTrustedIPs }}"
  - "--entrypoints.web.forwardedHeaders.trustedIPs=127.0.0.1/32,{{ .ForwardedHeadersTrustedIPs }}"
  {{- end }}
  {{- end }}
  {{- if .ProxyProtocol }}
  {{- if .ProxyProtocolInsecure }}
  - "--entrypoints.websecure.proxyProtocol.insecure"
  {{- end }}
  {{- if .ProxyProtocolTrustedIPs }}
  - "--entrypoints.websecure.proxyProtocol.trustedIPs=127.0.0.1/32,{{ .ProxyProtocolTrustedIPs }}"
  {{- end }}
  {{- end }}
  {{- if IngressProvider }}
  - "{{ printf "%s=%s" "--providers.kubernetesIngress.ingressClass" .IngressProvider }}"
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
