// Package git
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package git

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/pkg/git"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

var (
	gitRepoUrl        string
	gitClonePath      string
	gitAskCredentials bool
	gitUsername       string
	gitPassword       string
)

// installChart represents the list command
var clone = &cobra.Command{
	Use:   "clone",
	Short: "clone repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloning repository: ", gitRepoUrl)

		if gitRepoUrl == "" {
			panic("Repository URL is required")
		}

		if gitClonePath == "" {
			panic("Clone Path is required")
		}

		if gitAskCredentials {
			username, password, err := credentials()
			if err != nil {
				fmt.Println(err)
			}
			gitUsername = username
			gitPassword = password
		}

		if gitUsername != "" && gitPassword != "" {
			err := git.CloneWithCredentials(gitRepoUrl, gitClonePath, gitUsername, gitPassword)
			if err != nil {
				panic(err)
			}
			return
		}

		err := git.Clone(gitRepoUrl, gitClonePath)
		if err != nil {
			panic(err)
		}
	},
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func init() {
	clone.Flags().StringVarP(&gitRepoUrl, "url", "u", "", "Repository URL")
	clone.Flags().StringVarP(&gitClonePath, "path", "p", "", "Clone Path")
	clone.Flags().BoolVarP(&gitAskCredentials, "ask-credentials", "", false, "Ask for credentials")
	clone.Flags().StringVarP(&gitUsername, "username", "", "", "Username")
	clone.Flags().StringVarP(&gitPassword, "password", "", "", "Password")

	Cmd.AddCommand(clone)
}
