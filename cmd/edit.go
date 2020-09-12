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
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:       "edit [current|last|<time-entry-id>]",
	Aliases:   []string{"update"},
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"last", "current"},
	Short:     `Edit a time entry, use id "current" to apply to time entry in progress`,
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error
		interactive := viper.GetBool("interactive")

		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		tei, err := getTimeEntry(
			args[0],
			viper.GetString("workspace"),
			userID,
			c,
		)

		if err != nil {
			return err
		}

		if cmd.Flags().Changed("project") {
			tei.ProjectID, _ = cmd.Flags().GetString("project")
			if viper.GetBool("allow-project-name") && tei.ProjectID != "" {
				tei.ProjectID, err = getProjectByNameOrId(c, tei.WorkspaceID, tei.ProjectID)
				if err != nil && !interactive {
					return err
				}
			}
		}

		if cmd.Flags().Changed("description") {
			tei.Description, _ = cmd.Flags().GetString("description")
		}

		if cmd.Flags().Changed("task") {
			tei.TaskID, _ = cmd.Flags().GetString("task")
		}

		if cmd.Flags().Changed("tag") {
			tei.TagIDs, _ = cmd.Flags().GetStringSlice("tag")
		}

		if cmd.Flags().Changed("not-billable") {
			b, _ := cmd.Flags().GetBool("not-billable")
			tei.Billable = !b
		}

		if cmd.Flags().Changed("when") {
			whenString, _ = cmd.Flags().GetString("when")
			var v time.Time
			if v, err = convertToTime(whenString); err != nil {
				return err
			}
			tei.TimeInterval.Start = v
		}

		if cmd.Flags().Changed("end-at") {
			whenString, _ = cmd.Flags().GetString("end-at")
			var v time.Time
			if v, err = convertToTime(whenString); err != nil {
				return err
			}
			tei.TimeInterval.End = &v
		}

		if interactive {
			tei, err = confirmEntryInteractively(c, tei)
			if err != nil {
				return err
			}
		}

		tei, err = c.UpdateTimeEntry(api.UpdateTimeEntryParam{
			Workspace:   tei.WorkspaceID,
			TimeEntryID: tei.ID,
			Description: tei.Description,
			Start:       tei.TimeInterval.Start,
			End:         tei.TimeInterval.End,
			Billable:    tei.Billable,
			ProjectID:   tei.ProjectID,
			TaskID:      tei.TaskID,
			TagIDs:      tei.TagIDs,
		})

		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		return printTimeEntryImpl(c, tei, asJSON, format)
	}),
}

func init() {
	rootCmd.AddCommand(editCmd)

	addTimeEntryFlags(editCmd)

	editCmd.Flags().StringP("project", "p", "", "change the project")
	_ = completion.AddSuggestionsToFlag(editCmd, "project", suggestWithClientAPI(suggestProjects))

	editCmd.Flags().String("description", "", "change the description")
	editCmd.Flags().String("end-at", "", "when the entry should end (if not set \"\" will be used)")

	editCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	editCmd.Flags().BoolP("json", "j", false, "print as json")
}
