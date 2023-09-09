package haproxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/pkg/git"
	xos "github.com/zcubbs/zrun/pkg/os"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultLastCommitFile = "last_commit"
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

		style.PrintColoredHeader("configure haproxy from Git")

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

	// get user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %v", err)
	}

	// build path to lastCommitFile
	path := filepath.Join(homeDir,
		fmt.Sprintf("%s/%s", ".zrun", lastCommitFile),
	)

	// Create last commit file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := xos.CreateFileWithPath(path)
		if err != nil {
			return fmt.Errorf("failed to create last commit file: %v", err)
		}
	}

	lastCommit, err := os.ReadFile(lastCommitFile)
	if err != nil {
		fmt.Printf("warn: failed to read last commit: %v\n", err)
	}

	currentCommit, err := git.GetLatestCommit(clonePath)
	if err != nil {
		return fmt.Errorf("failed to get latest commit: %v", err)
	}

	fmt.Printf("last commit: %s\n", string(lastCommit))
	fmt.Printf("current commit: %s\n", currentCommit)

	if string(lastCommit) == "" {
		fmt.Println("First run")
	} else {
		changes, err := git.FileHasChanges(
			clonePath,
			file,
			strings.TrimSpace(string(lastCommit)),
			currentCommit,
		)
		if err != nil {
			return fmt.Errorf("failed to check if file has changes: %v", err)
		}

		if changes {
			fmt.Println("File has changes")
		} else {
			fmt.Println("File has no changes")
			return nil
		}
	}

	// Run the update action
	err = runAction(verbose)
	if err != nil {
		return fmt.Errorf("failed to run haproxy cfg copy/restart: %v", err)
	}

	// Save the latest commit for the next check
	err = os.WriteFile(lastCommitFile, []byte(currentCommit), 0644)
	if err != nil {
		return fmt.Errorf("failed to write last commit file: %v", err)
	}

	// Clean cloned files
	err = os.RemoveAll(clonePath)
	if err != nil {
		return fmt.Errorf("failed to clean cloned files: %v", err)
	}

	return nil
}

func runAction(verbose bool) error {
	// copy file to /etc/haproxy/haproxy.cfg
	filePath := fmt.Sprintf("%s/%s", clonePath, file)
	err := xos.CopyFileToDestination(filePath, "/etc/haproxy/haproxy.cfg")
	if err != nil {
		return fmt.Errorf("failed to copy file to destination: %v", err)
	}

	// validate config
	err = validateHaproxyConfig(verbose)
	if err != nil {
		return fmt.Errorf("failed to validate haproxy config: %v", err)
	}

	// restart haproxy
	err = xos.RestartSystemdService("haproxy", verbose)
	if err != nil {
		return fmt.Errorf("failed to restart haproxy: %v", err)
	}

	return nil
}

func init() {
	genTempPath := os.TempDir() + "/tmp-git-clone-" + time.Now().Format("20060102150405")
	updateFromGitCmd.Flags().StringVarP(&repoUrl, "repo-url", "r", "", "git repo url")
	updateFromGitCmd.Flags().StringVarP(&file, "file", "f", "haproxy.cfg", "file to watch for changes")
	updateFromGitCmd.Flags().StringVarP(&credentialsUsername, "credentials-username", "u", "", "git username")
	updateFromGitCmd.Flags().StringVarP(&credentialsPassword, "credentials-password", "p", "", "git password")
	updateFromGitCmd.Flags().StringVarP(&clonePath, "clone-path", "c", genTempPath, "clone path")
	updateFromGitCmd.Flags().StringVarP(&lastCommitFile, "last-commit-file", "l", defaultLastCommitFile, "last commit file")

	// Required flags
	_ = updateFromGitCmd.MarkFlagRequired("repo-url")
	_ = updateFromGitCmd.MarkFlagRequired("file")

	Cmd.AddCommand(updateFromGitCmd)
}
