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

	"gopkg.in/yaml.v2"
)

var (
	secretKey   string // secret key
	secretValue string // secret value
)

// Add represents the os command
var Add = &cobra.Command{
	Use:   "add",
	Short: "Add vault secret",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := addSecret()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func addSecret() error {
	sv, err := initVault()
	if err != nil {
		return err
	}

	err = sv.AddSecret(secretKey, secretValue)
	if err != nil {
		return err
	}

	err = sv.Save()
	if err != nil {
		return err
	}

	return nil
}

func initVault() (*zvault.SecretVault, error) {
	cfg, err := readVaultConfig()
	if err != nil {
		fmt.Println(err)
	}

	if cfg != nil {
		// set environment variables
		err = os.Setenv("VAULT_FILE", cfg.VaultFile)
		if err != nil {
			return nil, err
		}
		err = os.Setenv("VAULT_KEY", cfg.VaultKey)
		if err != nil {
			return nil, err
		}
	}

	filename := os.Getenv("VAULT_FILE")
	key := os.Getenv("VAULT_KEY")
	if filename == "" || key == "" {
		fmt.Println("creating new vault")
		filename = fmt.Sprintf("%s/%s", getUserHomePath(), ".zrun_vault")

		sv, err := zvault.InitializeVaultWithRandomKey(filename)
		if err != nil {
			return nil, err
		}

		err = writeVaultConfig(os.Getenv("VAULT_KEY"),
			filename,
			fmt.Sprintf("%s/%s", getUserHomePath(),
				".zrun_vault_cfg"),
		)

		return sv, nil
	}

	fmt.Println("using existing vault")
	sv, err := zvault.NewSecretVault()
	if err != nil {
		return nil, err
	}

	return sv, nil
}

func readVaultConfig() (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", getUserHomePath(),
		".zrun_vault_cfg"))
	if err != nil {
		return nil, err
	}

	// Create an empty Config instance
	cfg := &Config{}

	// Unmarshal the YAML file into the Config instance
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Config represents the YAML configuration file
type Config struct {
	VaultKey  string `yaml:"vault_key"`
	VaultFile string `yaml:"vault_file"`
}

func writeVaultConfig(keyStr, filename, yamlFile string) error {
	// Create a Config instance and populate it
	cfg := &Config{
		VaultKey:  keyStr,
		VaultFile: filename,
	}

	// Marshal the Config instance to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Write the data to the YAML file
	err = os.WriteFile(yamlFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getUserHomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func init() {
	// add flags
	Add.Flags().StringVar(&secretKey, "key", "", "secret key")
	Add.Flags().StringVar(&secretValue, "val", "", "secret value")

	// add mandatory flags
	if err := Add.MarkFlagRequired("key"); err != nil {
		fmt.Println(err)
	}

	if err := Add.MarkFlagRequired("val"); err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(Add)
}
