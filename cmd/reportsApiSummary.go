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
	"io"
	"os"
	"time"

	"github.com/lucassabreu/clockify-cli/reportsapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportsApiSummaryCmd represents the summary command
var reportsApiSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "A brief description of your command",
	RunE: withClockifyReportsClient(func(cmd *cobra.Command, args []string, c *reportsapi.ReportsClient) error {
		s, err := c.Summary(reportsapi.SummaryParam{
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

func init() {
	reportsApiCmd.AddCommand(reportsApiSummaryCmd)
}
