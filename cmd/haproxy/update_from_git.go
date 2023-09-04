package haproxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/pkg/git"
	xos "github.com/zcubbs/zrun/pkg/os"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"os"
	"strings"
	"time"
)

const (
	defaultLastCommitFile = "~/.zrun_haproxy_last_commit"
)

var (
	repoUrl             string // git repo url
	file                string // specific file to use
	credentialsUsername string // git username
	credentialsPassword string // git password
	clonePath           string // clone path
	lastCommitFile      string // last commit file
)

// updateFromGitCmd represents the list command
var updateFromGitCmd = &cobra.Command{
	Use:   "update-from-git",
	Short: "update haproxy config from git",
	Long:  `update haproxy config from git`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("configure haproxy for k3s")

		util.Must(
			util.RunTask(func() error {
				err := updateConfig(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func updateConfig(verbose bool) error {
	err := git.CloneWithCredentials(repoUrl, clonePath, credentialsUsername, credentialsPassword)
	if err != nil {
		return err
	}

	lastCommit, err := os.ReadFile(lastCommitFile)
	if err != nil {
		fmt.Printf("Failed to read last commit: %v\n", err)
		// Handle this based on your needs. If file doesn't exist, it's probably the first run.
	}

	currentCommit, err := git.GetLatestCommit(clonePath)
	if err != nil {
		return err
	}

	changes, err := git.FileHasChanges(
		clonePath,
		file,
		strings.TrimSpace(string(lastCommit)),
		currentCommit,
	)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if changes {
		err := runAction(verbose)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No changes found")
	}

	// Save the latest commit for the next check
	err = os.WriteFile(lastCommitFile, []byte(currentCommit), 0644)
	if err != nil {
		return err
	}

	// Clean cloned files
	err = os.RemoveAll(clonePath)
	if err != nil {
		return err
	}

	return nil
}

func runAction(verbose bool) error {
	// copy file to /etc/haproxy/haproxy.cfg
	// validate config
	// restart haproxy

	err := xos.CopyFileToDestination(file, "/etc/haproxy/haproxy.cfg")
	if err != nil {
		return err
	}

	err = validateHaproxyConfig(verbose)
	if err != nil {
		return err
	}

	err = xos.RestartSystemdService("haproxy", verbose)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	genTempPath := os.TempDir() + "/tmp-git-clone-" + time.Now().Format("20060102150405")
	updateFromGitCmd.Flags().StringVarP(&repoUrl, "repo-url", "r", "", "git repo url")
	updateFromGitCmd.Flags().StringVarP(&file, "file", "f", "", "file to watch for changes")
	updateFromGitCmd.Flags().StringVarP(&credentialsUsername, "credentials-username", "u", "", "git username")
	updateFromGitCmd.Flags().StringVarP(&credentialsPassword, "credentials-password", "p", "", "git password")
	updateFromGitCmd.Flags().StringVarP(&clonePath, "clone-path", "c", genTempPath, "clone path")
	updateFromGitCmd.Flags().StringVarP(&lastCommitFile, "last-commit-file", "l", defaultLastCommitFile, "last commit file")

	// Required flags
	_ = updateFromGitCmd.MarkFlagRequired("repo-url")
	_ = updateFromGitCmd.MarkFlagRequired("file")

	Cmd.AddCommand(updateFromGitCmd)
}
