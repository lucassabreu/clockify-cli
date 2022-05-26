package thismonth

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// reportThisMonthCmd represents the reports this-month command
func NewCmdThisMonth(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "this-month",
		Short: "List all time entries in this month",
		RunE: func(cmd *cobra.Command, args []string) error {
			first, last := timehlp.GetMonthRange(timehlp.Today())
			return util.ReportWithRange(f, first, last, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
