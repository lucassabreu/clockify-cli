// Copyright © 2019 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/output"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportCmd represents the output command
var reportCmd = &cobra.Command{
	Use:   "report <start> <end>",
	Short: "List all time entries in the date ranges and with more data (format date as 2016-01-02)",
	Args:  cobra.ExactArgs(2),
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		start, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			return err
		}
		end, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			return err
		}

		return reportWithRange(c, start, end, cmd)
	}),
}

// reportThisMonthCmd represents the reports this-month command
var reportThisMonthCmd = &cobra.Command{
	Use:   "this-month",
	Short: "List all time entries in this month",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		first, last := getMonthRange(time.Now())
		return reportWithRange(c, first, last, cmd)
	}),
}

// reportLastMonthCmd represents the report last-month command
var reportLastMonthCmd = &cobra.Command{
	Use:   "last-month",
	Short: "List all time entries in last month",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		first, last := getMonthRange(time.Now().AddDate(0, -1, 0))
		return reportWithRange(c, first, last, cmd)
	}),
}

// reportThisWeekCmd represents the report last-month command
var reportThisWeekCmd = &cobra.Command{
	Use:   "this-week",
	Short: "List all time entries in this week",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		first, last := getWeekRange(time.Now())
		return reportWithRange(c, first, last, cmd)
	}),
}

// reportLastWeekCmd represents the report last-month command
var reportLastWeekCmd = &cobra.Command{
	Use:   "last-week",
	Short: "List all time entries in last week",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		first, last := getWeekRange(time.Now().AddDate(0, 0, -7))
		return reportWithRange(c, first, last, cmd)
	}),
}

// reportLastDayCmd represents the report last-day command
var reportLastDayCmd = &cobra.Command{
	Use:   "last-day",
	Short: "List time entries from last day were a time entry exists",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		u, err := getUserId(c)
		if err != nil {
			return err
		}
		te, err := getTimeEntry(
			"last",
			viper.GetString(WORKSPACE),
			u,
			c,
		)
		if err != nil {
			return err
		}

		return reportWithRange(c, te.TimeInterval.Start, te.TimeInterval.Start, cmd)
	}),
}

// reportLastWeekDayCmd represents the report last working week day command
var reportLastWeekDayCmd = &cobra.Command{
	Use:   "last-week-day",
	Short: "List time entries from last week day (use `clockify-cli config workweek-days` command to set then)",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		workweek := strhlp.Map(strings.ToLower, viper.GetStringSlice(WORKWEEK_DAYS))
		if len(workweek) == 0 {
			return errors.New("no workweek days were set")
		}

		day := truncateDate(time.Now()).Add(-1)
		if strhlp.Search(strings.ToLower(day.Weekday().String()), workweek) != -1 {
			return reportWithRange(c, day, day, cmd)
		}

		dayWeekday := int(day.Weekday())
		if dayWeekday == int(time.Sunday) {
			dayWeekday = int(time.Saturday + 1)
		}

		lastWeekDay := int(time.Sunday)
		for _, w := range workweek {
			if i := strhlp.Search(w, weekdays); i > lastWeekDay && i < dayWeekday {
				lastWeekDay = i
			}
		}

		day = day.Add(time.Duration(-24*(dayWeekday-lastWeekDay)) * time.Hour)
		return reportWithRange(c, day, day, cmd)
	}),
}

func init() {
	rootCmd.AddCommand(reportCmd)

	_ = reportCmd.MarkFlagRequired(WORKSPACE)
	_ = reportCmd.MarkFlagRequired(USER_ID_FLAG)

	reportFlags(reportCmd)

	reportCmd.AddCommand(reportFlags(reportThisMonthCmd))
	reportCmd.AddCommand(reportFlags(reportLastMonthCmd))
	reportCmd.AddCommand(reportFlags(reportThisWeekCmd))
	reportCmd.AddCommand(reportFlags(reportLastWeekCmd))
	reportCmd.AddCommand(reportFlags(reportLastDayCmd))
	reportCmd.AddCommand(reportFlags(reportLastWeekDayCmd))
}

func reportFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("fill-missing-dates", "e", false, "add empty lines for dates without time entries")

	return cmd
}

func getMonthRange(ref time.Time) (first time.Time, last time.Time) {
	first = ref.AddDate(0, 0, ref.Day()*-1+1)
	last = first.AddDate(0, 1, -1)

	return
}

func getWeekRange(ref time.Time) (first time.Time, last time.Time) {
	first = ref.AddDate(0, 0, int(ref.Weekday())*-1)
	last = first.AddDate(0, 0, 7)

	return
}

func reportWithRange(c *api.Client, start, end time.Time, cmd *cobra.Command) error {
	format, _ := cmd.Flags().GetString("format")
	asJSON, _ := cmd.Flags().GetBool("json")
	asCSV, _ := cmd.Flags().GetBool("csv")
	fillMissingDates, _ := cmd.Flags().GetBool("fill-missing-dates")

	start = truncateDate(start)
	end = truncateDate(end).Add(time.Hour * 24)

	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	log, err := c.LogRange(api.LogRangeParam{
		Workspace:       viper.GetString(WORKSPACE),
		UserID:          userId,
		FirstDate:       start,
		LastDate:        end,
		PaginationParam: api.PaginationParam{AllPages: true},
	})

	if err != nil {
		return err
	}

	sort.Slice(log, func(i, j int) bool {
		return log[j].TimeInterval.Start.After(
			log[i].TimeInterval.Start,
		)
	})

	if fillMissingDates && len(log) > 0 {
		newLog := make([]dto.TimeEntry, 0, len(log))

		fillMissing := func(user *dto.User, first, last time.Time) []dto.TimeEntry {
			first = time.Date(first.Year(), first.Month(), first.Day(), 0, 0, 0, 0, time.Local)
			last = time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, time.Local)
			d := int(last.Sub(first).Hours() / 24)
			if d <= 0 {
				return []dto.TimeEntry{}
			}

			missing := make([]dto.TimeEntry, d)
			for i, t := range missing {
				ti := first.AddDate(0, 0, i)
				t.TimeInterval.Start = ti
				t.TimeInterval.End = &ti
				missing[i] = t
			}
			return missing
		}

		nextDay := start
		for _, t := range log {
			newLog = append(newLog, fillMissing(t.User, nextDay, t.TimeInterval.Start)...)
			newLog = append(newLog, t)
			nextDay = t.TimeInterval.Start.Add(time.Duration(24-t.TimeInterval.Start.Hour()) * time.Hour)
		}
		log = append(newLog, fillMissing(log[0].User, nextDay, end)...)
	}

	var fn func([]dto.TimeEntry, io.Writer) error
	fn = output.TimeEntriesPrintWithTimeFormat(output.TIME_FORMAT_FULL)
	if asJSON {
		fn = output.TimeEntriesJSONPrint
	}

	if asCSV {
		fn = output.TimeEntriesCSVPrint
	}

	if format != "" {
		fn = output.TimeEntriesPrintWithTemplate(format)
	}

	return fn(log, os.Stdout)
}

func truncateDate(t time.Time) time.Time {
	return t.Truncate(time.Hour * 24)
}
