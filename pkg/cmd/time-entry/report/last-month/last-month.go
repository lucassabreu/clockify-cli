package lastmonth

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdLastMonth represents the report last-month command
func NewCmdLastMonth(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "last-month",
		Short: "List all time entries in last month",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			first, last := timehlp.GetMonthRange(timehlp.Today().AddDate(0, -1, 0))
			return util.ReportWithRange(f, first, last, cmd, of)
		},
	}

	util.AddReportFlags(f, cmd, &of)

	return cmd

}
