// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"github.com/zcubbs/zrun/defaults"
	"os/exec"
)

func Install(packages ...string) error {
	for _, p := range packages {
		stdout, err := exec.Command(
			defaults.BinSh,
			"-c",
			fmt.Sprintf("sudo apt install -y %s", p)).Output()
		if err != nil {
			return err
		}
		fmt.Println(string(stdout))
	}
	return nil
}

func Update() error {
	stdout, err := exec.Command(defaults.BinSh, "-c", "sudo apt update -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

func Upgrade() error {
	stdout, err := exec.Command(defaults.BinSh, "-c", "sudo apt upgrade -y").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))
	return nil
}
