// Package helm
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/configs"
	"github.com/zcubbs/zrun/helm"
	zTime "github.com/zcubbs/zrun/time"
	"log"
	"os"
	"time"
)

// list represents the list command
var list = &cobra.Command{
	Use:   "list",
	Short: "List all helm releases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteHelmListCmd()

	},
}

func ExecuteHelmListCmd() {
	_releases, err := helm.GetAllReleases(
		configs.Config.Kubeconfig.Path,
	)

	if err != nil {
		log.Fatal("Could not list helm releases", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Namespace", "Release", "Name", "Version", "Status", "App Version", "Deployed"})
	for _, release := range _releases {
		t.AppendRows([]table.Row{
			{
				release.Namespace,
				release.Name,
				release.Chart.Metadata.Name,
				release.Chart.Metadata.Version,
				release.Info.Status,
				release.Chart.Metadata.AppVersion,
				zTime.TimeElapsed(time.Now(), release.Info.FirstDeployed.Time, false),
			},
		})

	}
	t.AppendSeparator()
	t.SetStyle(table.StyleColoredDark)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateRows = false

	t.Render()
	if err != nil {
		println(err.Error())
		panic("Could not list helm releases")
	}
}

func init() {
	Cmd.AddCommand(list)
}
