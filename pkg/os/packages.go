// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"os/exec"
)

const (
	BinSH = "/bin/sh"
)

func Install(packages ...string) error {
	for _, p := range packages {
		stdout, err := exec.Command(BinSH, "-c",
			fmt.Sprintf("apt install -y %s", p)).Output()
		if err != nil {
			return err
		}
		fmt.Println(string(stdout))
	}
	return nil
}

func Update() error {
	stdout, err := exec.Command(BinSH, "-c", "apt update -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

func Upgrade() error {
	stdout, err := exec.Command(BinSH, "-c", "apt upgrade -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}
