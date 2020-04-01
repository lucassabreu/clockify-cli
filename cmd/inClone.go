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
	"strings"

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
		var err error

		workspace := viper.GetString("workspace")
		tec, err := getTimeEntry(
			args[0],
			workspace,
			viper.GetString("user.id"),
			c,
		)
		tec.TimeInterval.End = nil

		if err != nil {
			printError(err)
			return
		}

		if tec.TimeInterval.Start, err = convertToTime(whenString); err != nil {
			printError(err)
			return
		}

		tec, err = newEntry(c, tec, viper.GetBool("interactive"), !viper.GetBool("no-closing"))

		if err != nil {
			printError(err)
			return
		}

		te, err := c.ConvertIntoFullTimeEntry(tec)
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

func getTimeEntry(id, workspace, userID string, c *api.Client) (dto.TimeEntryImpl, error) {
	id = strings.ToLower(id)

	if id != "last" {
		tei, err := c.GetTimeEntry(api.GetTimeEntryParam{
			Workspace:   workspace,
			TimeEntryId: id,
		})

		if err != nil {
			return dto.TimeEntryImpl{}, err
		}

		if tei == nil {
			return dto.TimeEntryImpl{}, errors.New("no previous time entry found")
		}

		return *tei, nil
	}

	list, err := c.GetRecentTimeEntries(api.GetRecentTimeEntries{
		Workspace:    workspace,
		UserID:       userID,
		Page:         1,
		ItemsPerPage: 1,
	})

	if err != nil {
		return dto.TimeEntryImpl{}, err
	}

	if len(list.TimeEntriesList) == 0 {
		return dto.TimeEntryImpl{}, errors.New("there is no previous time entry")
	}

	return list.TimeEntriesList[0], err
}

func init() {
	rootCmd.AddCommand(inCloneCmd)
	inCmd.AddCommand(inCloneCmd)

	inCloneCmd.Flags().String("when", "", "when the entry should be closed, if not informed will use current time")

	inCloneCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	inCloneCmd.Flags().BoolP("json", "j", false, "print as json")
}
