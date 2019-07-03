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

// inCloneCmd represents the clone command
var inCloneCmd = &cobra.Command{
	Use:   "clone <time-entry-id>",
	Short: "Copy a time entry and starts it (use \"last\" to copy the last one)",
	Args:  cobra.ExactArgs(1),
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		var whenDate *time.Time
		var err error

		if whenDate, err = getDateTimeParam("Start", true, whenString, convertToTime); err != nil {
			printError(err)
			return
		}

		workspace := viper.GetString("workspace")
		tec, err := getTimeEntry(
			args[0],
			workspace,
			viper.GetString("user.id"),
			c,
		)

		if err != nil {
			printError(err)
			return
		}

		err = c.Out(api.OutParam{
			Workspace: workspace,
			End:       *whenDate,
		})

		if err != nil {
			printError(errors.New("can not end current time entry"))
			return
		}

		tei, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   workspace,
			Billable:    tec.Billable,
			Start:       *whenDate,
			ProjectID:   tec.ProjectID,
			Description: tec.Description,
			TagIDs:      tec.TagIDs,
			TaskID:      tec.TaskID,
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
	rootCmd.AddCommand(inCloneCmd)
	inCmd.AddCommand(inCloneCmd)
	inCloneCmd.Flags().String("when", "", "when the entry should be closed, if not informed will use current time")

	inCloneCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	inCloneCmd.Flags().BoolP("json", "j", false, "print as json")
}
