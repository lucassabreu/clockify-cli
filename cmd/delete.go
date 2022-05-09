// Copyright © 2020 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
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

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:       "delete [" + ALIAS_CURRENT + "|<time-entry-id>]...",
	Aliases:   []string{"del", "rm", "remove"},
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{ALIAS_CURRENT},
	Short:     `Delete time entry(ies), use id "` + ALIAS_CURRENT + `" to apply to time entry in progress`,
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		for i := range args {
			param := api.DeleteTimeEntryParam{
				Workspace:   viper.GetString(WORKSPACE),
				TimeEntryID: args[i],
			}

			if param.TimeEntryID == ALIAS_CURRENT {
				te, err := c.GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
					Workspace: param.Workspace,
					UserID:    viper.GetString(USER_ID),
				})

				if err != nil {
					return err
				}

				if te == nil {
					return errors.New("there is no time entry in progress")
				}

				param.TimeEntryID = te.ID
			}

			if err := c.DeleteTimeEntry(param); err != nil {
				return err
			}
		}
		return nil
	}),
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
