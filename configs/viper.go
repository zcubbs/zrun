// Package configs
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var cfgFile string

var Config Configuration

var (
	defaults = map[string]interface{}{
		"debug.enabled":   true,
		"kubeconfig.path": getUserHomePath() + "/.kube/config",
	}
	envPrefix   = "Z"
	configName  = "config"
	configType  = "yaml"
	configPaths = []string{
		".",
		fmt.Sprintf("%s/.zrun", getUserHomePath()),
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
			log.Println("Error loading .env file")
		}
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		for _, p := range configPaths {
			viper.AddConfigPath(p)
		}
		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
		err := viper.ReadInConfig()
		if err != nil {
			if viper.GetString("debug.enabled") == "true" {
				fmt.Println(err)
			}
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
		log.Fatalf("could not decode config into struct: %v", err)
	}
}
