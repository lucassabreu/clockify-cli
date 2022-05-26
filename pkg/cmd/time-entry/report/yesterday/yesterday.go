package yesterday

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdYesterday represents report today command
func NewCmdYesterday(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "List all time entries created yesterday",
		RunE: func(cmd *cobra.Command, args []string) error {
			day := timehlp.TruncateDate(timehlp.Today()).Add(-1)
			return util.ReportWithRange(f, day, day, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
