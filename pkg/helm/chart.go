// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/

package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	helmValues "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"os"
)

type InstallChartOptions struct {
	Kubeconfig   string
	RepoName     string
	RepoUrl      string
	ChartName    string
	Namespace    string
	ChartVersion string
	ChartValues  helmValues.Options
	Debug        bool
	Upgrade      bool
}

func InstallChart(options InstallChartOptions) error {
	var settings = cli.New()
	settings.KubeConfig = options.Kubeconfig
	settings.Debug = options.Debug

	actionConfig := new(action.Configuration)

	helmLog := chooseLogFunc(options.Debug)

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		options.Namespace,
		os.Getenv("HELM_DRIVER"),
		helmLog); err != nil {
		return fmt.Errorf("failed to init action config. %s", err)
	}

	if options.Upgrade {
		return updateChart(options, actionConfig, settings)
	}

	return installChart(options, actionConfig, settings)
}

func updateChart(options InstallChartOptions, actionConfig *action.Configuration, settings *cli.EnvSettings) error {
	client := action.NewUpgrade(actionConfig)
	client.Namespace = options.Namespace

	// If version is specified, set it
	if options.ChartVersion != "" {
		client.Version = options.ChartVersion
	}

	// Locate the chart
	chartPath, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", options.RepoName, options.ChartName), settings)
	if err != nil {
		return fmt.Errorf("failed to locate chart: %w", err)
	}

	// Load the chart
	loadedChart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart from path: %w", err)
	}

	// Merge provided values
	p := getter.All(settings)
	vals, err := options.ChartValues.MergeValues(p)
	if err != nil {
		return fmt.Errorf("failed to merge values: %w", err)
	}

	// Update the chart
	_, err = client.Run(options.ChartName, loadedChart, vals)
	if err != nil {
		return fmt.Errorf("failed to update chart: %w", err)
	}

	return nil
}

func installChart(options InstallChartOptions, actionConfig *action.Configuration, settings *cli.EnvSettings) error {
	client := action.NewInstall(actionConfig)
	client.CreateNamespace = true

	setClientOptions(client, options)

	chartPath, err := locateChartPath(client, options, settings)
	if err != nil {
		return err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart. %s", err)
	}

	if valid, err := isChartInstallable(chartRequested); !valid {
		return fmt.Errorf("chart %s is not installable. %s", chartRequested.Name(), err)
	}

	err = handleDependencies(client, settings, chartRequested, chartPath, getter.All(settings))
	if err != nil {
		return err
	}

	vals, err := options.ChartValues.MergeValues(getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to merge values: %w", err)
	}

	r, err := client.Run(chartRequested, vals)
	if err != nil {
		return handleInstallError(err, options)
	}

	printDebugInfo(r, options)

	return nil
}

func setClientOptions(client *action.Install, options InstallChartOptions) {
	client.IsUpgrade = options.Upgrade
	client.ReleaseName = options.ChartName
	client.Namespace = options.Namespace

	if options.ChartVersion != "" {
		client.Version = options.ChartVersion
	}

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
}

func locateChartPath(client *action.Install, options InstallChartOptions, settings *cli.EnvSettings) (string, error) {
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", options.RepoName, options.ChartName), settings)
	if err != nil {
		return "", fmt.Errorf("failed to locate chart. %s", err)
	}
	if options.Debug {
		helmLog := chooseLogFunc(options.Debug)
		helmLog("CHART PATH: %s\n", cp)
	}
	return cp, nil
}

func handleInstallError(err error, options InstallChartOptions) error {
	if err.Error() == "cannot re-use a name that is still in use" && options.Debug {
		fmt.Println("warning: chart release name already exists, no action taken!")
		return nil
	}
	return fmt.Errorf("failed to install chart: %w", err)
}

func printDebugInfo(release *release.Release, options InstallChartOptions) {
	if options.Debug {
		fmt.Println(release.Manifest)
		fmt.Println(release.Info)
		fmt.Println(release.Chart.Values)
	}
}

// isChartInstallable checks if the chart is installable.
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, fmt.Errorf("invalid chart type: %q", ch.Metadata.Type)
}

func handleDependencies(client *action.Install,
	settings *cli.EnvSettings,
	ch *chart.Chart,
	cp string,
	p getter.Providers) error {

	reqs := ch.Metadata.Dependencies
	if reqs == nil {
		// No dependencies to handle
		return nil
	}

	if err := action.CheckDependencies(ch, reqs); err != nil {
		if client.DependencyUpdate {
			return updateDependencies(cp, client, settings, p)
		} else {
			return fmt.Errorf("chart dependencies are not satisfied: %w", err)
		}
	}

	return nil
}

func updateDependencies(cp string, client *action.Install, settings *cli.EnvSettings, p getter.Providers) error {
	man := &downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        cp,
		Keyring:          client.ChartPathOptions.Keyring,
		SkipUpdate:       false,
		Getters:          p,
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}
	if err := man.Update(); err != nil {
		return fmt.Errorf("failed to update the chart dependencies: %w", err)
	}
	return nil
}

func UninstallChart(kubeconfig, name, namespace string, debug bool) error {
	var settings = cli.New()
	settings.KubeConfig = kubeconfig
	actionConfig := new(action.Configuration)

	var helmLog = noLog

	if debug {
		helmLog = debugLog
	}

	if err := actionConfig.Init(settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		helmLog); err != nil {
		return fmt.Errorf("failed to initialize helm action configuration: %w", err)
	}
	client := action.NewUninstall(actionConfig)

	release, err := client.Run(name)
	if err != nil {
		return fmt.Errorf("failed to uninstall chart: %w", err)
	}

	if debug {
		fmt.Println("uninstalled", release.Release.Name)
	}
	return nil
}
