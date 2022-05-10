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
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:       "edit [" + ALIAS_CURRENT + "|" + ALIAS_LAST + "|<time-entry-id>]",
	Aliases:   []string{"update"},
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{ALIAS_LAST, ALIAS_CURRENT},
	Short:     `Edit a time entry, use id "` + ALIAS_CURRENT + `" to apply to time entry in progress`,
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error

		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		tei, err := getTimeEntry(
			args[0],
			viper.GetString(WORKSPACE),
			userID,
			c,
		)

		if err != nil {
			return err
		}

		dc := newDescriptionCompleter(c, tei.WorkspaceID, tei.UserID)

		if tei, err = manageEntry(
			tei,
			fillTimeEntryWithFlags(cmd.Flags()),
			getPropsInteractiveFn(c, dc),
			getDatesInteractiveFn(),
			getAllowNameForIDsFn(c),
			getValidateTimeEntryFn(c),
		); err != nil {
			return err
		}

		if tei, err = c.UpdateTimeEntry(api.UpdateTimeEntryParam{
			Workspace:   tei.WorkspaceID,
			TimeEntryID: tei.ID,
			Description: tei.Description,
			Start:       tei.TimeInterval.Start,
			End:         tei.TimeInterval.End,
			Billable:    tei.Billable,
			ProjectID:   tei.ProjectID,
			TaskID:      tei.TaskID,
			TagIDs:      tei.TagIDs,
		}); err != nil {
			return err
		}

		return printTimeEntryImpl(tei, c, cmd, output.TIME_FORMAT_SIMPLE)
	}),
}

func init() {
	rootCmd.AddCommand(editCmd)

	addTimeEntryFlags(editCmd)

	editCmd.Flags().StringP("when", "s", "", "when the entry should be started")
	editCmd.Flags().StringP("when-to-close", "e", "", "when the entry should be closed")

	editCmd.Flags().String("end-at", "", `when the entry should end (if not set "" will be used)`)
	_ = editCmd.Flags().MarkDeprecated("end-at", "use `when-to-close` flag instead")
}
