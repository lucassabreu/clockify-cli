package report

import (
	"time"

	"github.com/MakeNowJust/heredoc"
	lastday "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/last-day"
	lastmonth "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/last-month"
	lastweek "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/last-week"
	lastweekday "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/last-week-day"
	thismonth "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/this-month"
	thisweek "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/this-week"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/today"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/yesterday"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdReport represents the reports command
func NewCmdReport(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "report [<start>] [<end>]",
		Short: "List all time entries in the date ranges and with more " +
			"data (format date as 2016-01-02)",
		Args:    cobra.MaximumNArgs(2),
		Aliases: []string{"log"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			start := time.Now()
			if len(args) > 0 {
				start, err = time.Parse("2006-01-02", args[0])
				if err != nil {
					return err
				}
			}

			end := start
			if len(args) > 1 {
				if args[1] == "now" || args[1] == "today" {
					end = time.Now()
				} else if args[1] == "yesterday" {
					end = timehlp.TruncateDate(time.Now()).Add(-1)
				} else if end, err = time.Parse(
					"2006-01-02", args[1]); err != nil {
					return err
				}
			}

			return util.ReportWithRange(f, start, end, cmd)
		},
	}

	cmd.Long = cmd.Short + "\n" + heredoc.Doc(`
		If no parameter is set, shows today's time entries
		Aliases today/now can be used for <end> argument to represent current date
		Alias yesterday can be used for <end> argument to represent previous date
	`)

	cmd.AddCommand(thismonth.NewCmdThisMonth(f))
	cmd.AddCommand(lastmonth.NewCmdLastMonth(f))
	cmd.AddCommand(thisweek.NewCmdThisWeek(f))
	cmd.AddCommand(lastweek.NewCmdLastWeek(f))
	cmd.AddCommand(lastday.NewCmdLastDay(f))
	cmd.AddCommand(lastweekday.NewCmdLastWeekDay(f))
	cmd.AddCommand(today.NewCmdToday(f))
	cmd.AddCommand(yesterday.NewCmdYesterday(f))

	util.AddReportFlags(f, cmd)
	_ = cmd.MarkFlagRequired("workspace")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
