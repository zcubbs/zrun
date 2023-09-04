// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"os"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyFileToDestination(srcFile, destFile string) error {
	input, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(destFile, input, 0644)
	if err != nil {
		return err
	}
	return nil
}
