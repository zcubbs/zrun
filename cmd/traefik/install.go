// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/cmd/helm"
	"github.com/zcubbs/zrun/configs"
	helmPkg "github.com/zcubbs/zrun/helm"
	"helm.sh/helm/v3/pkg/cli/values"
	"log"
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
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install traefik Chart",
	Long:  `install traefik Chart. Note: requires helm`,
	Run: func(cmd *cobra.Command, args []string) {
		err := installChart()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func installChart() error {
	fmt.Println("installing cert-manager Chart")
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
		options.Values = append(options.Values, "logs.access.enabled=false")
		options.Values = append(options.Values, "ingressRoute.dashboard.enabled=true")
		options.Values = append(options.Values, "persistence.enabled=false")
		options.Values = append(options.Values, "service.type=LoadBalancer")
		options.Values = append(options.Values, "service.enabled=true")

		additionalArgs = append(additionalArgs, defaultArgs[:]...)
	}

	// check if ingressProvider is set, if so, use it
	if ingressProvider != "" {
		additionalArgs = append(additionalArgs, fmt.Sprintf(
			"%s=%s",
			"--providers.kubernetesIngress.ingressClass",
			ingressProvider,
		))
	}

	options.Values = append(options.Values, addAdditionalArgs(additionalArgs)...)

	verbose := Cmd.Flag("verbose").Value.String() == "true"

	helm.ExecuteInstallChartCmd(helmPkg.InstallChartOptions{
		Kubeconfig:   kubeconfig,
		RepoName:     "traefik",
		RepoUrl:      "https://helm.traefik.io/traefik",
		ChartName:    "traefik",
		Namespace:    "traefik",
		ChartVersion: chartVersion,
		ChartValues:  options,
		Debug:        verbose,
	})

	return nil
}

func addAdditionalArgs(additionalArgs []string) []string {
	var args []string
	for i, arg := range additionalArgs {
		adaptedArg := fmt.Sprintf("%s[%d]=%s", "additionalArguments", i, arg)
		args = append(args, adaptedArg)
	}

	fmt.Println(args)
	return args
}

func init() {
	// parse flags
	install.Flags().StringVar(&chartVersion, "version", "", "chart version")
	install.Flags().StringArrayVar(&options.Values, "set", nil, "chart values")
	install.Flags().StringArrayVar(&additionalArgs, "set-arg", nil, "chart values additional arguments")
	install.Flags().BoolVar(&useDefaults, "defaults", false, "use default values")
	install.Flags().BoolVar(&withInsecure, "insecure", false, "use insecure connection")
	install.Flags().BoolVar(&withForwardedHeaders, "forwardedHeaders", false, "use insecure forwarded headers")
	install.Flags().BoolVar(&withProxyProtocol, "proxy", false, "use proxy protocol")
	install.Flags().StringVar(&ingressProvider, "ingressProvider", "", "ingress provider")

	Cmd.AddCommand(install)
}