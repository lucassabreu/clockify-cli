package thismonth

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// reportThisMonthCmd represents the reports this-month command
func NewCmdThisMonth(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "this-month",
		Short: "List all time entries in this month",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			first, last := timehlp.GetMonthRange(timehlp.Today())
			return util.ReportWithRange(f, first, last, cmd, of)
		},
	}

	cmd.Long = cmd.Short + "\n\n" +
		util.HelpNamesForIds + "\n" +
		util.HelpMoreInfoAboutPrinting

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
