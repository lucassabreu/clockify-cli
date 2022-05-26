package thisweek

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdThisWeek represents the report this-week command
func NewCmdThisWeek(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "this-week",
		Short: "List all time entries in this week",
		RunE: func(cmd *cobra.Command, args []string) error {
			first, last := timehlp.GetWeekRange(timehlp.Today())
			return util.ReportWithRange(f, first, last, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
