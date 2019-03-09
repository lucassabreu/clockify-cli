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

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logInProgressCmd represents the logInProgress command
var logInProgressCmd = &cobra.Command{
	Use:   "in-progress",
	Short: "Show time entry in progress (if any)",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		tei, err := c.LogInProgress(api.LogInProgressParam{
			Workspace: viper.GetString("workspace"),
		})

		if err != nil {
			printError(err)
			return
		}

		var te *dto.TimeEntry
		if tei != nil {
			tef, err := c.ConvertIntoFullTimeEntry(*tei)
			if err != nil {
				printError(err)
				return
			}
			te = &tef
		}

		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")

		var reportFn func(*dto.TimeEntry, io.Writer) error

		reportFn = reports.TimeEntryPrint

		if asJSON {
			reportFn = reports.TimeEntryJSONPrint
		}

		if format != "" {
			reportFn = reports.TimeEntryPrintWithTemplate(format)
		}

		if err = reportFn(te, os.Stdout); err != nil {
			printError(err)
		}
	}),
}

func init() {
	logCmd.AddCommand(logInProgressCmd)
	logInProgressCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	logInProgressCmd.Flags().BoolP("json", "j", false, "print as json")

	logInProgressCmd.MarkFlagRequired("workspace")
}
