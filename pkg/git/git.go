// Package git
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/zcubbs/zrun/pkg/bash"
	"os"
)

// Clone clones a git repository to a given path
func Clone(url string, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return err
}

func CloneWithCredentials(url string, path string, username string, password string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	fmt.Println("Cloned repository to " + path + " successfully!")

	return err
}

func GetLatestCommit(gitRepoPath string) (string, error) {
	commit, err := bash.ExecuteCmdWithOutput("git", "-C", gitRepoPath, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}

	return commit, nil
}
