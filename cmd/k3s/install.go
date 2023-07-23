// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/k3s"
	"log"
)

var (
	disable             []string
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
		err := k3s.Install(k3s.Config{
			Disable:                 disable,
			DataDir:                 dataDir,
			DefaultLocalStoragePath: volumeStorageDir,
			WriteKubeconfigMode:     writeKubeconfigMode,
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	install.Flags().StringArrayVarP(&disable, "disable", "dis", nil, "list of k3s features to disable")
	install.Flags().StringVarP(&dataDir, "dataDir", "dir", "", "data storage directory")
	install.Flags().StringVarP(&volumeStorageDir, "volumeStorageDir", "vs", "", "volume storage directory")
	install.Flags().StringVarP(&writeKubeconfigMode, "writeKubeconfigMode", "wm", "", "write kubeconfig mode")
	Cmd.AddCommand(install)
}
