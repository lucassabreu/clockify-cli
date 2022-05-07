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
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editMultipleCmd represents the editMultiple command
var editMultipleCmd = &cobra.Command{
	Use:       "edit-multiple [current|last|<time-entry-id>]...",
	Aliases:   []string{"update-multiple", "multi-edit", "multi-update", "mult-edit", "mult-update"},
	Args:      cobra.MinimumNArgs(2),
	ValidArgs: []string{"last", "current"},
	Short:     `Edit multiple time entries at once, use id "current"/"last" to apply to time entry in progress.`,
	Long: `edit multiple time entries at once, use id "current"/"last" to apply to time entry in progress.
When multiple IDs are informed the default values on interactive mode will be the values of the first time entry informed.
When using interactive mode all entries will end with the same properties except for Start and End, if you wanna edit only some properties, than use the flags without interactive mode.
Start and end fields can't be mass-edited.`,
	PreRunE: printMultipleTimeEntriesPreRun,
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error

		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		teis := make([]dto.TimeEntryImpl, len(args))
		for i := range args {
			tei, err := getTimeEntry(
				args[i],
				viper.GetString(WORKSPACE),
				userID,
				false,
				c,
			)
			if err != nil {
				return err
			}
			teis[i] = tei
		}

		tei := teis[0]
		editFn := func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			return c.UpdateTimeEntry(api.UpdateTimeEntryParam{
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
		}

		fn := func(input dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			var err error
			for i, tei := range teis {
				input.TimeInterval = tei.TimeInterval
				input.ID = tei.ID

				if tei, err = editFn(input); err != nil {
					return input, err
				}

				teis[i] = tei
			}

			return input, err
		}

		tei, err = fillTimeEntryWithFlags(tei, cmd.Flags())
		if err != nil {
			return err
		}

		interactive := viper.GetBool(INTERACTIVE)
		if !interactive {
			fn = func(input dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
				for i, tei := range teis {
					if cmd.Flags().Changed("project") {
						tei.ProjectID = input.ProjectID
					}

					if cmd.Flags().Changed("description") {
						tei.Description = input.Description
					}

					if cmd.Flags().Changed("task") {
						tei.TaskID = input.TaskID
					}

					if cmd.Flags().Changed("tag") || cmd.Flags().Changed("tags") {
						tei.TagIDs = input.TagIDs
					}

					if cmd.Flags().Changed("not-billable") {
						tei.Billable = input.Billable
					}

					teis[i] = tei
					if _, err = editFn(tei); err != nil {
						return tei, err
					}
				}
				return input, nil
			}
		}

		dc := newDescriptionCompleter(c, tei.WorkspaceID, tei.UserID)

		if _, err = manageEntry(
			tei,
			fn,
			getInteractiveFn(c, dc, false),
			getAllowNameForIDsFn(c),
			getValidateTimeEntryFn(c),
		); err != nil {
			return err
		}

		tes := make([]dto.TimeEntry, len(teis))
		var t *dto.TimeEntry
		for i, tei := range teis {
			t, err = c.GetHydratedTimeEntry(api.GetTimeEntryParam{
				TimeEntryID: tei.ID,
				Workspace:   tei.WorkspaceID,
			})

			if err != nil {
				return err
			}
			tes[i] = *t
		}

		return printTimeEntries(tes, cmd, output.TIME_FORMAT_SIMPLE)
	}),
}

func init() {
	rootCmd.AddCommand(editMultipleCmd)

	addTimeEntryFlags(editMultipleCmd)
	addPrintMultipleTimeEntriesFlags(editMultipleCmd)
}
