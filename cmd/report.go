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
	"os"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportsCmd represents the reports command
var reportCmd = &cobra.Command{
	Use:   "report <start> <end>",
	Short: "report for date ranges and with more data (format date as 2016-01-02)",
	Args:  cobra.ExactArgs(2),
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		start, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			printError(err)
			return
		}
		end, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			printError(err)
			return
		}

		reportWithRange(c, start, end)
	}),
}

// reportThisMonthCmd represents the reports this-month command
var reportThisMonthCmd = &cobra.Command{
	Use:   "this-month",
	Short: "report all entries in this month",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		first, last := getMonthRange(time.Now())
		reportWithRange(c, first, last)
	}),
}

// reportLastMonthCmd represents the report last-month command
var reportLastMonthCmd = &cobra.Command{
	Use:   "last-month",
	Short: "report all entries in last month",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		first, last := getMonthRange(time.Now().AddDate(0, -1, 0))
		reportWithRange(c, first, last)
	}),
}

func init() {
	rootCmd.AddCommand(reportCmd)

	_ = reportCmd.MarkFlagRequired("workspace")
	_ = reportCmd.MarkFlagRequired("user-id")

	reportCmd.AddCommand(reportThisMonthCmd)
	reportCmd.AddCommand(reportLastMonthCmd)
}

func getMonthRange(ref time.Time) (first time.Time, last time.Time) {
	first = ref.AddDate(0, 0, ref.Day()*-1+1).Truncate(time.Hour)
	last = first.AddDate(0, 1, -1)

	return
}

func reportWithRange(c *api.Client, start, end time.Time) {
	log, err := c.LogRange(api.LogRangeParam{
		Workspace: viper.GetString("workspace"),
		UserID:    viper.GetString("user.id"),
		FirstDate: start,
		LastDate:  end,
		AllPages:  true,
	})

	if err != nil {
		printError(err)
		return
	}

	reports.TimeEntriesPrint(log, os.Stdout)
}
