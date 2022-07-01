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
	of := util.NewOutputFlags()
	cmd := &cobra.Command{
		Use:   "report [<start>] [<end>]",
		Short: "List all time entries for a given date range",
		Long: heredoc.Docf(`
			List all time entries for a given date range

			If no parameter is set, shows today's time entries
			Aliases today/now can be used for <end> argument to represent current date
			Alias yesterday can be used for <end> argument to represent previous date

			To choose a specific date to start or end use the format "2006-01-02"

			%s
			All the subcommands have the same flags to filter and format the time entries, but will act as aliases to relative date ranges.
		`, util.HelpNamesForIds),
		Example: heredoc.Docf(`
			# reporting all time entries from today
			$ %[1]s
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			|            ID            |        START        |         END         |   DUR   |   PROJECT    |          DESCRIPTION           |              TAGS              |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b87a9785815e619d7ce02e | 2022-06-26 12:25:56 | 2022-06-26 12:26:47 | 0:00:51 | Clockify Cli | Example for today              | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b87abb85815e619d7ce034 | 2022-06-26 12:26:47 | 2022-06-26 13:00:00 | 0:33:13 | Clockify Cli | Example for today (second one) | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| TOTAL                    |                     |                     | 0:34:04 |              |                                |                                |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+

			# reporting all time entries from 2022-06-24 to today
			$ %[1]s 2022-06-24 today
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			|            ID            |        START        |         END         |   DUR   |   PROJECT    |          DESCRIPTION           |              TAGS              |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b8ce7185815e619d7d0a82 | 2022-06-24 08:00:00 | 2022-06-24 09:00:00 | 1:00:00 | Clockify Cli | Example for before yesterday   | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b8ce1edba0da0f21e7e688 | 2022-06-25 08:00:00 | 2022-06-25 09:00:00 | 1:00:00 | Clockify Cli | Example for yesterday          | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b87a9785815e619d7ce02e | 2022-06-26 12:25:56 | 2022-06-26 12:26:47 | 0:00:51 | Clockify Cli | Example for today              | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| 62b87abb85815e619d7ce034 | 2022-06-26 12:26:47 | 2022-06-26 13:00:00 | 0:33:13 | Clockify Cli | Example for today (second one) | Development                    |
			|                          |                     |                     |         |              |                                | (62ae28b72518aa18da2acb49)     |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+
			| TOTAL                    |                     |                     | 2:34:04 |              |                                |                                |
			+--------------------------+---------------------+---------------------+---------+--------------+--------------------------------+--------------------------------+

			# when there are no entries for the range
			$ %[1]s 1999-01-01
			+-------+-------+-----+---------+---------+-------------+------+
			|  ID   | START | END |   DUR   | PROJECT | DESCRIPTION | TAGS |
			+-------+-------+-----+---------+---------+-------------+------+
			| TOTAL |       |     | 0:00:00 |         |             |      |
			+-------+-------+-----+---------+---------+-------------+------+

			# format output with golang template
			$ %[1]s 2022-06-23 --format "{{.ID}} - {{ .TimeInterval.Duration }} - {{ pad .Project.Name 12 }} - {{ .Description }}"
			62b8d162984dba2c06699e3f - PT1H - Clockify Cli - First example for report
			62b8d195dba0da0f21e7e85d - PT1H - Special      - Lunch break
			62b8d207dba0da0f21e7e868 - PT1H - Clockify Cli - After lunch

			# show time spent on the project "Clockify CLI" as float
			$ %[1]s 2022-06-23 --duration-float -p "clockify cli"
			2.000000

			# show time spent on the project "Clockify CLI" with "lunch" on description
			$ %[1]s 2022-06-23 --duration-formatted -p "clockify cli" -d lunch
			1:00:00

			# show ids from time entries from project "clockify cli"
			$ %[1]s 2022-06-23 -p "clockify cli" --quiet
			62b8d162984dba2c06699e3f
			62b8d207dba0da0f21e7e868

			# show time entries from project "special" as markdown
			$ %[1]s 2022-06-23 -p "clockify cli" --quiet
			ID: %[2]s63b8d195dba0da0f21e7e85d%[2]s  
			Billable: %[2]sno%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Special (%[2]s6202680228782767055ef004%[2]s)  
			Interval: %[2]s2023-06-23 15:00:00%[2]s until %[2]s2022-06-23 16:00:00%[2]s  
			Description:
			> Lunch break

			Tags:
			 * Meeting (%[2]s6219486e8cb9606d934ebb5f%[2]s)

			# csv format output
			$ %[1]s --csv
			id,description,project.id,project.name,task.id,task.name,start,end,duration,user.id,user.email,user.name,tags...
			62b87a9785815e619d7ce02e,Example for today,621948458cb9606d934ebb1c,Clockify Cli,62b87a7e984dba2c0669724d,Report Command,2022-06-26 12:25:56,2022-06-26 12:26:47,0:00:51,5c6bf21db079873a55facc08,joe@due.com,John Due,Development (62ae28b72518aa18da2acb49)
			62b87abb85815e619d7ce034,Example for today (second one),621948458cb9606d934ebb1c,Clockify Cli,62b87a7e984dba2c0669724d,Report Command,2022-06-26 12:26:47,2022-06-26 13:00:00,0:33:13,5c6bf21db079873a55facc08,joe@due.com,John Due,Development (62ae28b72518aa18da2acb49)
		`, "clockify-cli report", "`"),
		Args:    cobra.MaximumNArgs(2),
		Aliases: []string{"log"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			var err error

			start := timehlp.Today()
			if len(args) > 0 {
				start, err = time.Parse("2006-01-02", args[0])
				if err != nil {
					return err
				}
			}

			end := start
			if len(args) > 1 {
				if args[1] == "now" || args[1] == "today" {
					end = timehlp.Today()
				} else if args[1] == "yesterday" {
					end = timehlp.Today().Add(-1)
				} else if end, err = time.Parse(
					"2006-01-02", args[1]); err != nil {
					return err
				}
			}

			return util.ReportWithRange(f, start, end, cmd, of)
		},
	}

	cmd.AddCommand(thismonth.NewCmdThisMonth(f))
	cmd.AddCommand(lastmonth.NewCmdLastMonth(f))
	cmd.AddCommand(thisweek.NewCmdThisWeek(f))
	cmd.AddCommand(lastweek.NewCmdLastWeek(f))
	cmd.AddCommand(lastday.NewCmdLastDay(f))
	cmd.AddCommand(lastweekday.NewCmdLastWeekDay(f))
	cmd.AddCommand(today.NewCmdToday(f))
	cmd.AddCommand(yesterday.NewCmdYesterday(f))

	util.AddReportFlags(f, cmd, &of)
	_ = cmd.MarkFlagRequired("workspace")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
