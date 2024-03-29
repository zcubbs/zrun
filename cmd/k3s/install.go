// Package k3s
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package k3s

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/x/k3s"
	"github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
	"github.com/zcubbs/zrun/internal/configs"
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
		style.PrintColoredHeader("install k3s")
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		must.Succeed(
			progress.RunTask(func() error {
				err := k3s.Install(k3s.Config{
					Disable:                 disable,
					TlsSan:                  tlsSan,
					DataDir:                 dataDir,
					DefaultLocalStoragePath: volumeStorageDir,
					WriteKubeconfigMode:     writeKubeconfigMode,
				}, verbose)
				if err != nil {
					return err
				}

				kubeconfig := configs.Config.Kubeconfig.Path
				ok, err := kubernetes.IsClusterReady(cmd.Context(), kubeconfig)
				if !ok || err != nil {
					return err
				}
				return nil
			}, true))
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
