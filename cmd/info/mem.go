package info

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"os"
)

// memUsageCmd represents the memUsage command
var memUsageCmd = &cobra.Command{
	Use:   "mem",
	Short: "memory usage info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := mem.VirtualMemory()

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Total", "Free", "Used%"})
		t.AppendRows([]table.Row{
			{
				humanize.Bytes(v.Total),
				humanize.Bytes(v.Free),
				fmt.Sprintf("%.0f%%", v.UsedPercent),
			},
		})
		t.AppendSeparator()
		t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
		t.Render()
	},
}

func init() {
	Cmd.AddCommand(memUsageCmd)
}
