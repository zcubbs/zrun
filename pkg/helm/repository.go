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
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string, debug bool) error {
	var settings = cli.New()
	repoFile := settings.RepositoryConfig
	repoFile = filepath.Clean(repoFile)

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory %s. %s", filepath.Dir(repoFile), err)
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
				panic(fmt.Errorf("failed to unlock file %s. %s", repoFile, err))
			}
		}(fileLock)
	}
	if err != nil {
		return fmt.Errorf("failed to lock file %s. %s", repoFile, err)
	}

	b, err := os.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read file %s. %s", repoFile, err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("failed to unmarshal file %s. %s", repoFile, err)
	}

	if f.Has(name) {
		if debug {
			fmt.Printf("%q already exists with the same configuration, skipping\n", name)
		}
		return nil
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return fmt.Errorf("looks like %q is not a valid chart repository or cannot be reached. %s", url, err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return fmt.Errorf("failed to write file %s. %s", repoFile, err)
	}

	if debug {
		fmt.Printf("%q has been added to your repositories\n", name)
	}

	return nil
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate(debug bool) error {
	var settings = cli.New()
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached",
				cfg.URL)
		}
		repos = append(repos, r)
	}

	if debug {
		fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	}

	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if debug {
				if _, err := re.DownloadIndexFile(); err != nil {
					fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
				} else {
					fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
				}
			}
		}(re)
	}

	wg.Wait()

	if debug {
		fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
	}

	return nil
}
