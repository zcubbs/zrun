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
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"helm.sh/helm/v3/pkg/cli/values"
)

const (
	traefikNamespace = "traefik"
)

var (
	defaultArgs = [...]string{
		"--global.sendanonymoususage=false",
		"--entrypoints.websecure.http.tls",
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
	chartVersion         string
	options              values.Options
	additionalArgs       []string
	useDefaults          bool
	withInsecure         bool
	withForwardedHeaders bool
	withProxyProtocol    bool
	ingressProvider      string
	endpointWeb          string
	endpointWebsecure    string
)

var (
	dnsProviderString string
	dnsResolver       string
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
		wa := fmt.Sprintf("%s=%s", "ports.web.exposedPort", endpointWeb)
		wsa := fmt.Sprintf("%s=%s", "ports.websecure.exposedPort", endpointWebsecure)
		options.Values = append(options.Values, "logs.access.enabled=false")
		options.Values = append(options.Values, "ingressRoute.dashboard.enabled=true")
		options.Values = append(options.Values, "persistence.enabled=false")
		options.Values = append(options.Values, "service.type=LoadBalancer")
		options.Values = append(options.Values, "service.enabled=true")
		options.Values = append(options.Values, wa)
		options.Values = append(options.Values, wsa)

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
		args, err := configureDNSChallenge()
		if err != nil {
			return err
		}

		options.Values = append(options.Values, "fromEnv[0].name=traefik-dns-account-credentials")
		additionalArgs = append(additionalArgs, args[:]...)

		err = createDnsSecret()
		if err != nil {
			return fmt.Errorf("failed to create dns secret: %w", err)
		}
	}

	args := addAdditionalArgs(additionalArgs)

	if verbose {
		fmt.Printf("...traefik additional args: %s\n", args)
	}

	options.Values = append(options.Values, args...)

	err := helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
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

func addAdditionalArgs(additionalArgs []string) []string {
	var args []string
	for i, arg := range additionalArgs {
		adaptedArg := fmt.Sprintf("%s[%d]=%s", "additionalArguments", i, arg)
		args = append(args, adaptedArg)
	}

	return args
}

func init() {
	// parse flags
	installCmd.Flags().StringVar(&chartVersion, "version", "", "chart version")
	installCmd.Flags().StringSliceVar(&options.Values, "set", nil, "chart values")
	installCmd.Flags().StringSliceVar(&additionalArgs, "set-arg", nil, "chart values additional arguments")
	installCmd.Flags().BoolVar(&useDefaults, "defaults", false, "use default values")
	installCmd.Flags().BoolVar(&withInsecure, "insecure", false, "use insecure connection")
	installCmd.Flags().BoolVar(&withForwardedHeaders, "forwardedHeaders", false, "use insecure forwarded headers")
	installCmd.Flags().BoolVar(&withProxyProtocol, "proxy", false, "use proxy protocol")
	installCmd.Flags().StringVar(&ingressProvider, "ingressProvider", "", "ingress provider")
	installCmd.Flags().StringVar(&endpointWeb, "endpointWeb", "80", "endpoint web")
	installCmd.Flags().StringVar(&endpointWebsecure, "endpointWebsecure", "443", "endpoint websecure")

	installCmd.PersistentFlags().BoolVar(&useVault, "use-vault", false, "use vault")
	installCmd.PersistentFlags().StringVar(&dnsProviderString, "dns-provider", "", "dns provider")
	installCmd.PersistentFlags().StringVar(&dnsResolver, "dns-resolver", "letsencrypt", "dns resolver")
	installCmd.PersistentFlags().StringVar(&dnsTz, "dns-tz", "Europe/Paris", "dns tz")

	Cmd.AddCommand(installCmd)
}
