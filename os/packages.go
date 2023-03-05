// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"github.com/zcubbs/zrun/defaults"
	"os"
	"os/exec"
)

func Install(packages ...string) error {
	for _, p := range packages {
		cmd := &exec.Cmd{
			// Path is the path of the command to run.
			Path: defaults.BinBash,
			// Args holds command line arguments, including the command itself as Args[0].
			Args:   []string{"sudo", "apt", "install", "-y", p},
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func Update() error {
	cmd := &exec.Cmd{
		// Path is the path of the command to run.
		Path: defaults.BinBash,
		// Args holds command line arguments, including the command itself as Args[0].
		Args:   []string{"sudo", "apt", "update", "-y"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func Upgrade() error {
	cmd := &exec.Cmd{
		// Path is the path of the command to run.
		Path: defaults.BinBash,
		// Args holds command line arguments, including the command itself as Args[0].
		Args:   []string{"sudo", "apt", "upgrade", "-y"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
