// Package configs
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package configs

type Configuration struct {
	Kubeconfig  `mapstructure:"kubeconfig" json:"kubeconfig"`
	VaultConfig `mapstructure:"vault" json:"vault"`
}

type VaultConfig struct {
	Path string `mapstructure:"path" json:"path"`
	Key  string `mapstructure:"key" json:"key"`
}

type Kubeconfig struct {
	Path string `mapstructure:"path" json:"path"`
}
