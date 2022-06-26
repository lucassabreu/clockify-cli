package lastweekday

import (
	"errors"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewCmdLastWeekDay represents the report last working week day command
func NewCmdLastWeekDay(f cmdutil.Factory) *cobra.Command {
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "last-week-day",
		Short: "List time entries from last week day",
		Long: heredoc.Docf(`
			List time entries from last week day

			For the CLI to know which days of the week you are expected to work, you will need to set them.
			This can be done using:
			$ clockify-cli config init

			Or more directly by running the set command as follows:
			$ clockify-cli config set workweek-days monday,tuesday,wednesday,thursday,friday

			%s
			%s
		`,
			util.HelpNamesForIds,
			util.HelpMoreInfoAboutPrinting,
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			workweek := f.Config().GetWorkWeekdays()
			if len(workweek) == 0 {
				return errors.New("no workweek days were set")
			}

			day := timehlp.Today().Add(-1)
			if strhlp.Search(
				strings.ToLower(day.Weekday().String()), workweek) != -1 {
				return util.ReportWithRange(f, day, day, cmd, of)
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
			return util.ReportWithRange(f, day, day, cmd, of)
		},
	}

	util.AddReportFlags(f, cmd, &of)

	return cmd
}
