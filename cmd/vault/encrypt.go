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
	encryptFile string
	encryptKey  string
)

// update represents the list command
var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt file or string",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if encryptKey == "" {
			k, err := getEncryptKey()
			if err != nil {
				panic(err)
			}

			encryptKey = k
		}
		if encryptFile != "" {
			err := zv.EncryptFile(encryptFile, encryptKey)
			if err != nil {
				panic(err)
			}
		}

		s := args[0]
		if s == "" {
			panic("No string to encrypt")
		}

		es := zv.Encrypt(s, encryptKey)
		fmt.Println(es)
	},
}

func getEncryptKey() (string, error) {
	fmt.Print("encryption key: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func init() {
	encrypt.Flags().StringVarP(&encryptFile, "file", "f", "", "File to encrypt")
	encrypt.Flags().StringVarP(&encryptKey, "key", "k", "", "Key to encrypt with")

	Cmd.AddCommand(encrypt)
}
