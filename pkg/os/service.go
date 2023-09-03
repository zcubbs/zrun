// Package os provides a set of functions to interact with the operating system.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RestartSystemdService(serviceName string) error {
	// Restart the service
	_, err := executeCmd("systemctl", "restart", serviceName)
	if err != nil {

		return fmt.Errorf("Failed to restart service: %w\n", err)
	}

	// Check the status of the service
	status, err := executeCmd("systemctl", "status", serviceName)
	if err != nil {
		return fmt.Errorf("Failed to get service status: %w\n", err)
	}

	fmt.Println("Service status:")
	fmt.Println(status)

	return nil
}

func executeCmd(command string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}
