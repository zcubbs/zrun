// Package vault
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package vault

import (
	"fmt"
	"github.com/spf13/cobra"
	zv "github.com/zcubbs/zrun/vault"
	"golang.org/x/term"
	"strings"
	"syscall"
)

var (
	decryptFile string
	decryptKey  string
)

// update represents the list command
var decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt file or string",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if decryptKey == "" {
			k, err := getDecryptKey()
			if err != nil {
				panic(err)
			}

			decryptKey = k
		}
		if encryptFile != "" {
			s, err := zv.DecryptFile(decryptFile, decryptKey)
			if err != nil {
				panic(err)
			}

			fmt.Println(s)
			return
		}

		s := args[0]
		if s == "" {
			panic("No string to decrypt")
		}

		ds := zv.Decrypt(s, decryptKey)
		fmt.Println(ds)
	},
}

func getDecryptKey() (string, error) {
	fmt.Print("encryption key: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func init() {
	encrypt.Flags().StringVarP(&decryptFile, "file", "f", "", "File to decrypt")
	encrypt.Flags().StringVarP(&decryptKey, "key", "k", "", "Key to decrypt with")

	Cmd.AddCommand(decrypt)
}
