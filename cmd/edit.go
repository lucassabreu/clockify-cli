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
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit <time-entry-id>",
	Aliases: []string{"update"},
	Args:    cobra.ExactArgs(1),
	Short:   `Edit a time entry, use id "current" to apply to time entry in progress`,
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		var err error
		param := api.UpdateTimeEntryParam{
			Workspace:   viper.GetString("workspace"),
			TimeEntryID: args[0],
		}

		if param.TimeEntryID == "current" {
			te, err := c.LogInProgress(api.LogInProgressParam{
				Workspace: param.Workspace,
			})

			if err != nil {
				printError(err)
				return
			}

			if te == nil {
				printError(errors.New("there is no time entry in progress"))
				return
			}

			param.TimeEntryID = te.ID
		}

		param.ProjectID, _ = cmd.Flags().GetString("project")
		param.Description, _ = cmd.Flags().GetString("description")
		param.TaskID, _ = cmd.Flags().GetString("task")
		param.TagIDs, _ = cmd.Flags().GetStringSlice("tag")

		b, _ := cmd.Flags().GetBool("not-billable")
		param.Billable = !b

		whenString, _ = cmd.Flags().GetString("when")
		var v time.Time
		if v, err = convertToTime(whenString); err != nil {
			printError(err)
			return
		}
		param.Start = v

		if cmd.Flags().Changed("end-at") {
			whenString, _ = cmd.Flags().GetString("end-at")
			var v time.Time
			if v, err = convertToTime(whenString); err != nil {
				printError(err)
				return
			}
			param.End = &v
		}

		tei, err := c.UpdateTimeEntry(param)

		if err != nil {
			printError(err)
			return
		}

		te, err := c.ConvertIntoFullTimeEntry(tei)
		if err != nil {
			printError(err)
			return
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

		if err = reportFn(&te, os.Stdout); err != nil {
			printError(err)
		}
	}),
}

func init() {
	rootCmd.AddCommand(editCmd)

	addTimeEntryFlags(editCmd)

	editCmd.Flags().StringP("project", "p", "", "change the project")
	editCmd.Flags().String("description", "", "change the description")
	editCmd.Flags().String("end-at", "", "when the entry should end (if not set \"\" will be used)")

	editCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	editCmd.Flags().BoolP("json", "j", false, "print as json")
}
