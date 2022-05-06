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

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:       "edit [current|last|<time-entry-id>]",
	Aliases:   []string{"update"},
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"last", "current"},
	Short:     `Edit a time entry, use id "current" to apply to time entry in progress`,
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
			false,
			c,
		)

		if err != nil {
			return err
		}

		tei, err = fillTimeEntryWithFlags(tei, cmd.Flags())
		if err != nil {
			return err
		}

		var dc *descriptionCompleter
		if viper.GetBool(DESCR_AUTOCOMP) {
			dc = newDescriptionCompleter(
				c,
				tei.WorkspaceID,
				tei.UserID,
				viper.GetInt(DESCR_AUTOCOMP_DAYS),
			)
		}

		return manageEntry(
			c,
			tei,
			func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
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
			},
			getInteractiveFn(c, dc, true),
			getAllowNameForIDsFn(c),
			printTimeEntryImpl(c, cmd, output.TIME_FORMAT_SIMPLE),
			getValidateTimeEntryFn(c),
		)
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
