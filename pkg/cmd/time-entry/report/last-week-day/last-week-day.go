package lastweekday

import (
	"errors"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewCmdLastWeekDay represents the report last working week day command
func NewCmdLastWeekDay(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "last-week-day",
		Short: "List time entries from last week day " +
			"(use `clockify-cli config workweek-days` command to set then)",
		RunE: func(cmd *cobra.Command, args []string) error {
			workweek := f.Config().GetWorkWeekdays()
			if len(workweek) == 0 {
				return errors.New("no workweek days were set")
			}

			day := timehlp.TruncateDate(time.Now()).Add(-1)
			if strhlp.Search(
				strings.ToLower(day.Weekday().String()), workweek) != -1 {
				return util.ReportWithRange(f, day, day, cmd)
			}

			dayWeekday := int(day.Weekday())
			if dayWeekday == int(time.Sunday) {
				dayWeekday = int(time.Saturday + 1)
			}

			lastWeekDay := int(time.Sunday)
			for _, w := range workweek {
				i := strhlp.Search(w, cmdutil.GetWeekdays())
				if i > lastWeekDay && i < dayWeekday {
					lastWeekDay = i
				}
			}

			day = day.Add(
				time.Duration(-24*(dayWeekday-lastWeekDay)) * time.Hour)
			return util.ReportWithRange(f, day, day, cmd)
		},
	}

	util.AddReportFlags(f, cmd)

	return cmd
}
