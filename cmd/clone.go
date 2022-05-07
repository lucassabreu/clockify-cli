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

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:       "clone [last|<time-entry-id>]",
	Short:     `Copy a time entry and starts it (use "last" to copy the last one)`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"last"},
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error

		userId, err := getUserId(c)
		if err != nil {
			return err
		}
		workspace := viper.GetString(WORKSPACE)
		tec, err := getTimeEntry(
			args[0],
			workspace,
			userId,
			true,
			c,
		)
		tec.UserID = viper.GetString(USER_ID)
		tec.TimeInterval.End = nil

		if err != nil {
			return err
		}

		if tec, err = fillTimeEntryWithFlags(tec, cmd.Flags()); err != nil {
			return err
		}

		noClosing, _ := cmd.Flags().GetBool("no-closing")

		dc := newDescriptionCompleter(c, tec.WorkspaceID, tec.UserID)

		if !noClosing {
			if err := validateClosingTimeEntry(
				c, tec.WorkspaceID, viper.GetString(USER_ID),
			); err != nil {
				return err
			}
		}

		if tec, err = manageEntry(
			tec,
			getAllowNameForIDsFn(c),
			getPropsInteractiveFn(c, dc),
			getDatesInteractiveFn(),
			getValidateTimeEntryFn(c),
		); err != nil {
			return err
		}

		if !noClosing {
			if err = out(tec, c); err != nil {
				return err
			}
		}

		if tec, err = createTimeEntry(tec, c); err != nil {
			return err
		}

		return printTimeEntryImpl(tec, c, cmd, output.TIME_FORMAT_SIMPLE)
	}),
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	addTimeEntryFlags(cloneCmd)
	addTimeEntryDateFlags(cloneCmd)

	cloneCmd.Flags().BoolP("no-closing", "", false, "don't close any active time entry")
}
