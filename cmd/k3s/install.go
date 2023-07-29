// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
	"github.com/zcubbs/zrun/util"
)

var (
	disable             []string
	tlsSan              []string
	dataDir             string
	volumeStorageDir    string
	writeKubeconfigMode string
)

// install represents the list command
var install = &cobra.Command{
	Use:   "install",
	Short: "install k3s",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("-------------------------------------------")
		fmt.Println("installing k3s...")
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		util.Must(k3s.Install(k3s.Config{
			Disable:                 disable,
			TlsSan:                  tlsSan,
			DataDir:                 dataDir,
			DefaultLocalStoragePath: volumeStorageDir,
			WriteKubeconfigMode:     writeKubeconfigMode,
		}, verbose))
	},
}

func init() {
	install.Flags().StringArrayVar(&tlsSan, "tls-san", nil, "list of k3s tls to add to certificate")
	install.Flags().StringArrayVar(&disable, "disable", nil, "list of k3s features to disable")
	install.Flags().StringVar(&dataDir, "data-dir", "", "data storage directory")
	install.Flags().StringVar(&volumeStorageDir, "volume-storage-dir", "", "volume storage directory")
	install.Flags().StringVar(&writeKubeconfigMode, "write-kubeconfig-Mode", "", "write kubeconfig mode")
	Cmd.AddCommand(install)
}
