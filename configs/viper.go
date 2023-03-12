// Package configs
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var cfgFile string

var Config Configuration

const defaultCfgFileName = ".zrun"

var (
	defaults = map[string]interface{}{
		"debug.enabled":   true,
		"kubeconfig.path": "/etc/rancher/k3s/k3s.yaml",
	}
	envPrefix   = "Z"
	configName  = defaultCfgFileName
	configType  = "yaml"
	configPaths = []string{
		".",
		getUserHomePath(),
	}
)

var allowedEnvVarKeys = []string{
	"awx.url",
	"awx.username",
	"awx.password",
	"k3s.version",
	"debug.enabled",
	"kubeconfig.path",
}

func getUserHomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func Bootstrap() {
	err := godotenv.Load(".env")

	if err != nil {
		if viper.GetString("debug.enabled") == "true" {
			fmt.Printf("[info] loading .env file")
		}
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		initConfigFile()
		for _, p := range configPaths {
			viper.AddConfigPath(p)
		}
		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("[warn] %s\n", err)
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(envPrefix)

	for _, key := range allowedEnvVarKeys {
		err := viper.BindEnv(key)
		if err != nil {
			fmt.Println(err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&Config)
	if err != nil {
		fmt.Printf("[warn] could not decode config into struct: %v", err)
	}
}

func initConfigFile() {
	// create cfg if not exists
	path := filepath.Join(getUserHomePath(), defaultCfgFileName)
	path = filepath.Clean(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
	}
}
