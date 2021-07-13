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
	"github.com/lucassabreu/clockify-cli/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dateString string
var yesterday bool
var dateFormat = "2006-01-02"

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"logs"},
	Short:   "List the entries from a specific day",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		var filterDate time.Time

		var err error
		if filterDate, err = time.Parse(dateFormat, dateString); err != nil {
			return err
		}

		if yesterday {
			filterDate = time.Now().Add(time.Hour * -24)
		}

		userId, err := getUserId(c)
		if err != nil {
			return err
		}

		log, err := c.Log(api.LogParam{
			Workspace:       viper.GetString(WORKSPACE),
			UserID:          userId,
			Date:            filterDate,
			PaginationParam: api.PaginationParam{AllPages: true},
		})

		sort.Slice(log, func(i, j int) bool {
			return log[j].TimeInterval.Start.After(
				log[i].TimeInterval.Start,
			)
		})

		if err != nil {
			return err
		}

		var reportFn func([]dto.TimeEntry, io.Writer) error
		reportFn = output.TimeEntriesPrint

		if asJSON {
			reportFn = output.TimeEntriesJSONPrint
		}

		if format != "" {
			reportFn = output.TimeEntriesPrintWithTemplate(format)
		}

		return reportFn(log, os.Stdout)
	}),
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().StringVarP(&dateString, "date", "d", time.Now().Format(dateFormat), "set the date to be logged in the format: YYYY-MM-DD")
	logCmd.Flags().BoolVarP(&yesterday, "yesterday", "y", false, "list the yesterday's entries")
	logCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	logCmd.Flags().BoolP("json", "j", false, "print as json")

	_ = logCmd.MarkFlagRequired(WORKSPACE)
	_ = logCmd.MarkFlagRequired(USER_ID)
}
