// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	helmValues "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
}

func InstallChart(options InstallChartOptions) {
	var settings = cli.New()
	settings.KubeConfig = options.Kubeconfig
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), options.Namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(
			fmt.Errorf("failed to init action config. %s", err))
	}

	client := action.NewInstall(actionConfig)

	client.CreateNamespace = true

	if options.ChartVersion != "" {
		client.Version = options.ChartVersion
	}

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = options.ChartName
	cp, err := client.ChartPathOptions.LocateChart(
		fmt.Sprintf("%s/%s",
			options.RepoName,
			options.ChartName),
		settings)
	if err != nil {
		log.Fatal(err)
	}

	if options.Debug {
		debug("CHART PATH: %s\n", cp)
	}

	p := getter.All(settings)

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Fatal(err)
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Fatal(err)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		checkDependencies(client, settings, chartRequested, req, cp, p)
	}
	vals, err := options.ChartValues.MergeValues(p)
	if err != nil {
		log.Fatal("Failed to merge values: ", err)
	}

	client.Namespace = options.Namespace
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Fatal("Failed to install chart: ", err)
	}

	if options.Debug {
		fmt.Println(release.Manifest)
		fmt.Println(release.Info)
		fmt.Println(release.Chart.Values)
	}
}

func checkDependencies(client *action.Install,
	settings *cli.EnvSettings,
	ch *chart.Chart,
	reqs []*chart.Dependency,
	cp string,
	p getter.Providers) {
	if err := action.CheckDependencies(ch, reqs); err != nil {
		if client.DependencyUpdate {
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
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string) {
	var settings = cli.New()
	repoFile := settings.RepositoryConfig
	repoFile = filepath.Clean(repoFile)

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func(fileLock *flock.Flock) {
			err := fileLock.Unlock()
			if err != nil {
				log.Fatal(err)
			}
		}(fileLock)
	}
	if err != nil {
		log.Fatal(err)
	}

	b, err := os.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(name) {
		fmt.Printf("repository name (%s) already exists\n", name)
		return
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		log.Fatal(err)
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q has been added to your repositories\n", name)
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate() {
	var settings = cli.New()
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		log.Fatal(errors.New("no repositories found. You must add one before updating"))
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
}

func UninstallChart(kubeconfig, name, namespace string) {
	var settings = cli.New()
	settings.KubeConfig = kubeconfig
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}
	client := action.NewUninstall(actionConfig)

	release, err := client.Run(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uninstalled", release.Release.Name)
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	err := log.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		return
	}
}
