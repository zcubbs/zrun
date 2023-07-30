// Package bash
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package bash

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteScript(script string, output bool, commands ...string) (bool, error) {
	var out bytes.Buffer
	cmd := &exec.Cmd{
		Path:   script,
		Args:   commands,
		Stdout: &out,
		Stderr: &out,
	}

	if output {
		fmt.Println("Executing command ", cmd)
		fmt.Println(out.String())
	}

	err := cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExecuteCmd(cmd string, output bool, args ...string) error {
	execute := exec.Command(cmd, args...)
	stdout, err := execute.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Print the output
	if output {
		fmt.Printf("Executing command %s\n", execute.String())
		fmt.Println(string(stdout))
	}

	return nil
}
