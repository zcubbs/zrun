package info

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// diskUsageCmd represents the diskUsage command
var diskUsageCmd = &cobra.Command{
	Use:   "disk",
	Short: "disk usage info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		defaultDirectory := "."

		if dir := viper.GetString("cmd.info.diskUsage.defaultDirectory"); dir != "" {
			defaultDirectory = dir
		}
		usage := du.NewDiskUsage(defaultDirectory)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Size", "Used", "Avail", "Free", "Use%"})
		t.AppendRows([]table.Row{
			{
				humanize.Bytes(usage.Size()),
				humanize.Bytes(usage.Used()),
				humanize.Bytes(usage.Available()),
				humanize.Bytes(usage.Free()),
				fmt.Sprintf("%.0f%%", usage.Usage()*100),
			},
		})
		t.AppendSeparator()
		t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
		t.Render()
	},
}

func init() {
	Cmd.AddCommand(diskUsageCmd)
}
