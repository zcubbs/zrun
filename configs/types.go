// Package configs
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package configs

type Configuration struct {
	Kubeconfig `mapstructure:"kubeconfig" json:"kubeconfig"`
}

type Kubeconfig struct {
	Path string `mapstructure:"path" json:"path"`
}
