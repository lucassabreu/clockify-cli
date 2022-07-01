package yesterday

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdYesterday represents report today command
func NewCmdYesterday(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "List all time entries created yesterday",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			day := timehlp.Today().Add(-1)
			return util.ReportWithRange(f, day, day, cmd, of)
		},
	}

	cmd.Long = cmd.Short + "\n\n" +
		util.HelpNamesForIds + "\n" +
		util.HelpMoreInfoAboutPrinting

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
