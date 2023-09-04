package os

import (
	"fmt"
	"os/exec"
	"strings"
)

// AddCronJob adds a new cron job to the current user's crontab.
// The cronJob parameter should be the full string of the cron job, e.g., "* * * * * /path/to/script.sh".
func AddCronJob(cronJob string) error {
	// Get the current crontab
	currentCron, err := exec.Command("crontab", "-l").CombinedOutput()
	if err != nil && !strings.Contains(string(currentCron), "no crontab for") {
		return err
	}

	// Check if the cron job already exists to avoid duplicates
	if strings.Contains(string(currentCron), cronJob) {
		return fmt.Errorf("cron job already exists")
	}

	// Append the new cron job
	updatedCron := string(currentCron) + cronJob + "\n"

	// Save the updated crontab
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(updatedCron)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
