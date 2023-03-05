// Package bash
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package bash

import (
	"fmt"
	"log"
	"os"
)

func Chmod(path string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}
	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File permissions before: %s\n", stats.Mode())
	err = os.Chmod(path, perm)
	if err != nil {
		log.Fatal(err)
	}

	stats, err = os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File permissions after: %s\n", stats.Mode())

	return nil
}

func RmDir(path string) error {
	return os.RemoveAll(path)
}

func Mkdir(path string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}
	return os.MkdirAll(path, perm)
}

func FileExists(path string) bool {
	// check if file exists
	// if it does, return true
	// if it doesn't, return false
	// TODO: implement check if file exists

	return false
}

func DirExists(path string) bool {
	// check if directory exists
	// if it does, return true
	// if it doesn't, return false
	// TODO: implement check if directory exists

	return false
}
