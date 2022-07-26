package today

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdToday represents report today command
func NewCmdToday(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "today",
		Short: "List all time entries created today",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			today := timehlp.Today()
			return util.ReportWithRange(f, today, today, cmd, of)
		},
	}

	cmd.Long = cmd.Short + "\n\n" +
		util.HelpNamesForIds + "\n" +
		util.HelpMoreInfoAboutPrinting

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
