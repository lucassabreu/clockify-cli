package lastweek

import (
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdLastWeek represents the report last-week command
func NewCmdLastWeek(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-week",
		Short: "List all time entries in last week",
		RunE: func(cmd *cobra.Command, args []string) error {
			first, last := timehlp.GetWeekRange(
				timehlp.TruncateDate(time.Now()).AddDate(0, 0, -7))
			return util.ReportWithRange(f, first, last, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
