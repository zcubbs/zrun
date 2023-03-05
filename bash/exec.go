package bash

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecuteScript(script string, commands ...string) (bool, error) {
	cmd := &exec.Cmd{
		Path:   script,
		Args:   commands,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println("Executing command ", cmd)

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

func ExecuteCmd(cmd string, args ...string) error {
	execute := exec.Command(cmd, args...)
	stdout, err := execute.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Print the output
	fmt.Println(string(stdout))

	return nil
}
