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
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/lucassabreu/clockify-cli/api/dto"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

var tags []string
var notBillable bool
var noClosing bool
var task string

var whenString string
var whenToCloseString string

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:   "in <project-id> <description>",
	Short: "Create a new time entry and starts it",
	Args:  cobra.MaximumNArgs(2),
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {

		var whenToCloseDate time.Time
		var err error

		tei := dto.TimeEntryImpl{
			WorkspaceID:  viper.GetString("workspace"),
			TagIDs:       tags,
			TimeInterval: dto.TimeInterval{},
		}

		if len(args) > 0 {
			tei.ProjectID = args[0]
		}

		if len(args) > 1 {
			tei.Description = args[1]
		}

		if whenString != "" {
			tei.TimeInterval.Start, err = convertToTime(whenString)
			if err != nil {
				printError(fmt.Errorf("Fail to convert when to start: %s", err.Error()))
				return
			}
		}

		if whenToCloseString != "" {
			whenToCloseDate, err = convertToTime(whenToCloseString)
			if err != nil {
				printError(fmt.Errorf("Fail to convert when to end: %s", err.Error()))
				return
			}
			tei.TimeInterval.End = &whenToCloseDate
		}

		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		err = newEntry(c, tei, viper.GetBool("interactive"), !noClosing, format, asJSON)
		if err != nil {
			printError(err)
		}
	}),
}

func init() {
	rootCmd.AddCommand(inCmd)

	addTimeEntryFlags(inCmd)

	inCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	inCmd.Flags().BoolP("json", "j", false, "print as json")
	inCmd.PersistentFlags().BoolVar(&noClosing, "no-closing", false, "don't close any active time entry")

	_ = viper.BindPFlag("no-closing", inCmd.PersistentFlags().Lookup("no-closing"))
}

func addTimeEntryFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&notBillable, "not-billable", "n", false, "is this time entry not billable")
	cmd.Flags().StringVar(&task, "task", "", "add a task to the entry")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "add tags to the entry")
	cmd.Flags().StringVar(&whenString, "when", time.Now().Format(fullTimeFormat), "when the entry should be closed, if not informed will use current time")
}
