package thisweek

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdThisWeek represents the report this-week command
func NewCmdThisWeek(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "this-week",
		Short: "List all time entries in this week",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			first, last := timehlp.GetWeekRange(timehlp.Today())
			return util.ReportWithRange(f, first, last, cmd, of)
		},
	}

	cmd.Long = cmd.Short + "\n\n" +
		util.HelpNamesForIds + "\n" +
		util.HelpMoreInfoAboutPrinting

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
