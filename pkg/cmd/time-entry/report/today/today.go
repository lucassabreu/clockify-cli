package today

import (
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdToday represents report today command
func NewCmdToday(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "List all time entries created today",
		RunE: func(cmd *cobra.Command, args []string) error {
			today := time.Now()
			return util.ReportWithRange(f, today, today, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
