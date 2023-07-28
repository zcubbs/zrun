// Package vault
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package vault

import (
	"fmt"
	"github.com/spf13/cobra"
	zvault "github.com/zcubbs/zrun/vault"
	"os"
)

var (
	getSecretKey string // secret key
)

// Get represents the os command
var Get = &cobra.Command{
	Use:   "get",
	Short: "Get vault secret",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		secretValue, err := GetSecret(getSecretKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(secretValue)
	},
}

func GetSecret(key string) (string, error) {
	cfg, err := readVaultConfig()
	if err != nil {
		return "", err
	}

	// set environment variables
	err = os.Setenv("VAULT_FILE", cfg.VaultFile)
	if err != nil {
		return "", err
	}
	err = os.Setenv("VAULT_KEY", cfg.VaultKey)
	if err != nil {
		return "", err
	}

	sv, err := zvault.NewSecretVault()
	if err != nil {
		return "", err
	}

	secretValue, err := sv.GetSecret(key)
	if err != nil {
		return "", err
	}

	return secretValue, nil
}

func init() {
	// add flags
	Get.Flags().StringVar(&getSecretKey, "key", "", "secret key")

	// add mandatory flags
	err := Get.MarkFlagRequired("key")
	if err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(Get)
}
