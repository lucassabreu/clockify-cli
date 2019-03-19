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

	"github.com/spf13/viper"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

var cardNumber int
var issueNumber int
var tags []string
var notBillable bool
var task string

var whenString string
var whenToCloseString string

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:     "in <project-name-or-id> <description>",
	Short:   "Create a new time entry and starts it",
	Example: `clockify-cli in --issue 13 "timesheet"`,
	Args:    cobra.MinimumNArgs(1),
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {

		var whenDate time.Time
		var whenToCloseDate *time.Time
		var err error

		project := args[0]

		var description string
		if len(args) > 1 {
			description = args[1]
		}

		if whenDate, err = convertToTime(whenString); err != nil {
			printError(err)
			return
		}

		if whenToCloseString != "" {
			var t time.Time
			if t, err = convertToTime(whenToCloseString); err != nil {
				printError(err)
				return
			}
			whenToCloseDate = &t
		}

		workspace := viper.GetString("workspace")

		err = c.Out(api.OutParam{
			Workspace: workspace,
			End:       whenDate,
		})

		if err != nil {
			printError(errors.New("can not end current time entry"))
			return
		}

		tei, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   workspace,
			Billable:    !notBillable,
			Start:       whenDate,
			End:         whenToCloseDate,
			ProjectID:   project,
			Description: description,
			TagIDs:      tags,
			TaskID:      task,
		})

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
	rootCmd.AddCommand(inCmd)

	addTimeEntryFlags(inCmd)

	inCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	inCmd.Flags().BoolP("json", "j", false, "print as json")
}

func addTimeEntryFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&notBillable, "not-billable", "n", false, "is this time entry not billable")
	cmd.Flags().IntVarP(&cardNumber, "card", "c", 0, "trello card number being started")
	cmd.Flags().IntVarP(&issueNumber, "issue", "i", 0, "issue number being started")
	cmd.Flags().StringVar(&task, "task", "", "add a task to the entry")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "add tags to the entry")
	cmd.Flags().StringVar(&whenString, "when", time.Now().Format(fullTimeFormat), "when the entry should be closed, if not informed will use current time")
}
