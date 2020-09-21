// Copyright © 2020 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/reportsapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportsApiSummaryCmd represents the summary command
var reportsApiSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "A brief description of your command",
	RunE: withClockifyReportsClient(func(cmd *cobra.Command, args []string, c *reportsapi.ReportsClient) error {

		sgs, err := cmd.Flags().GetStringSlice("groups")
		if err != nil {
			return err
		}

		gs, err := stringsToGroupSlice(sgs)
		if err != nil {
			return err
		}

		s, err := c.Summary(reportsapi.SummaryParam{
			Groups: gs,
			TimeEntryParam: reportsapi.TimeEntryParam{
				Workspace:      viper.GetString("workspace"),
				DateRangeStart: time.Now().AddDate(0, 0, -5),
				DateRangeEnd:   time.Now(),
			},
		})

		summaryJSONPrint(s, os.Stdout)

		return err
	}),
}

func summaryJSONPrint(s reportsapi.SummaryReport, w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

func stringsToGroupSlice(ss []string) (reportsapi.GroupSlice, error) {
	gs := make(reportsapi.GroupSlice, len(ss))
	for i, s := range ss {
		switch strings.TrimSpace(strings.ToLower(s)) {
		case strings.ToLower(string(reportsapi.Project)):
			gs[i] = reportsapi.Project
		case strings.ToLower(string(reportsapi.Client)):
			gs[i] = reportsapi.Client
		case strings.ToLower(string(reportsapi.Task)):
			gs[i] = reportsapi.Task
		case strings.ToLower(string(reportsapi.Tag)):
			gs[i] = reportsapi.Tag
		case strings.ToLower(string(reportsapi.Date)):
			gs[i] = reportsapi.Date
		case strings.ToLower(string(reportsapi.User)):
			gs[i] = reportsapi.User
		case strings.ToLower(string(reportsapi.UserGroup)):
			gs[i] = reportsapi.UserGroup
		case strings.ToLower(string(reportsapi.TimeEntry)):
			gs[i] = reportsapi.TimeEntry
		default:
			return gs, fmt.Errorf(`"%s" is not a valid group`, s)
		}

	}

	return gs, nil
}

func addGroupReportFlag(cmd *cobra.Command) error {
	cmd.Flags().StringSliceP("groups", "g", []string{}, "add groups to the report")
	return completion.AddFixedSuggestionsToFlag(cmd, "groups", completion.ValigsArgsSlide{
		string(reportsapi.Project),
		string(reportsapi.Client),
		string(reportsapi.Task),
		string(reportsapi.Tag),
		string(reportsapi.Date),
		string(reportsapi.User),
		string(reportsapi.UserGroup),
		string(reportsapi.TimeEntry),
	})
}

func init() {
	reportsApiCmd.AddCommand(reportsApiSummaryCmd)
	_ = addGroupReportFlag(reportsApiSummaryCmd)
}
