package lastday

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdLastDay represents the report last-day command
func NewCmdLastDay(f cmdutil.Factory) *cobra.Command {
	of := util.NewReportFlags()
	cmd := &cobra.Command{
		Use:   "last-day",
		Short: "List time entries from last day were a time entry was created",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			u, err := f.GetUserID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			te, err := timeentryhlp.GetLatestEntryEntry(c, w, u)
			if err != nil {
				return err
			}

			return util.ReportWithRange(
				f, te.TimeInterval.Start, te.TimeInterval.Start, cmd, of)
		},
	}

	cmd.Long = cmd.Short + "\n\n" +
		util.HelpNamesForIds + "\n" +
		util.HelpMoreInfoAboutPrinting

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
