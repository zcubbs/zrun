package haproxy

import (
	"fmt"
	"github.com/spf13/cobra"
	xos "github.com/zcubbs/zrun/pkg/os"
	"github.com/zcubbs/zrun/pkg/style"
	"github.com/zcubbs/zrun/pkg/util"
)

// setupGitCronCmd represents the list command
var setupGitCronCmd = &cobra.Command{
	Use:   "setup-git-cron",
	Short: "setup haproxy git update cron",
	Long:  `setup haproxy git update cron`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"

		style.PrintColoredHeader("configure haproxy for k3s")

		util.Must(
			util.RunTask(func() error {
				err := setupCron(verbose)
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

func setupCron(_ bool) error {
	const scriptPath = "/etc/haproxy/update_from_git.sh"
	// generate script
	err := xos.GenerateBashScript(scriptPath,
		`zrun haproxy update-from-git \
		--repo-url $HAPROXY_GIT_REPO \
		--file $HAPROXY_GIT_CONFIG_FILE \
		--credentials-username $HAPROXY_GIT_USERNAME \
		--credentials-password $HAPROXY_GIT_PASSWORD \
		`)
	if err != nil {
		return err
	}

	// add cron job
	err = xos.AddCronJob(fmt.Sprintf("* * * * * %s", scriptPath))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Cmd.AddCommand(setupGitCronCmd)
}
