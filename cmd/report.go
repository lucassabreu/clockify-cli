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
	"io"
	"os"
	"sort"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportsCmd represents the reports command
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

func init() {
	rootCmd.AddCommand(reportCmd)

	_ = reportCmd.MarkFlagRequired("workspace")
	_ = reportCmd.MarkFlagRequired("user-id")

	reportFlags(reportCmd)
	reportFlags(reportThisMonthCmd)
	reportFlags(reportLastMonthCmd)
	reportFlags(reportThisWeekCmd)
	reportFlags(reportLastWeekCmd)

	reportCmd.AddCommand(reportThisMonthCmd)
	reportCmd.AddCommand(reportLastMonthCmd)
	reportCmd.AddCommand(reportThisWeekCmd)
	reportCmd.AddCommand(reportLastWeekCmd)
}

func reportFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("fill-missing-dates", "e", false, "add empty lines for dates without time entries")
}

func getMonthRange(ref time.Time) (first time.Time, last time.Time) {
	first = ref.AddDate(0, 0, ref.Day()*-1+1).Truncate(time.Hour)
	last = first.AddDate(0, 1, -1)

	return
}

func getWeekRange(ref time.Time) (first time.Time, last time.Time) {
	first = ref.AddDate(0, 0, int(ref.Weekday())*-1).Truncate(time.Hour)
	last = first.AddDate(0, 0, 5)

	return
}

func reportWithRange(c *api.Client, start, end time.Time, cmd *cobra.Command) error {
	format, _ := cmd.Flags().GetString("format")
	asJSON, _ := cmd.Flags().GetBool("json")
	asCSV, _ := cmd.Flags().GetBool("csv")
	fillMissingDates, _ := cmd.Flags().GetBool("fill-missing-dates")

	start = start.Add(time.Duration(start.Hour()) * time.Hour * -1)
	end = end.Add(time.Duration(24-start.Hour()) * time.Hour * 1)

	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	log, err := c.LogRange(api.LogRangeParam{
		Workspace:       viper.GetString("workspace"),
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
	fn = reports.TimeEntriesPrintWithTimeFormat(reports.TIME_FORMAT_FULL)
	if asJSON {
		fn = reports.TimeEntriesJSONPrint
	}

	if asCSV {
		fn = reports.TimeEntriesCSVPrint
	}

	if format != "" {
		fn = reports.TimeEntriesPrintWithTemplate(format)
	}

	return fn(log, os.Stdout)
}
