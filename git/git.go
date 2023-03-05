// Package git
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package git

import (
	"github.com/go-git/go-git/v5"
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
