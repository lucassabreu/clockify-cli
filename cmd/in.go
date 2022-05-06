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
	"github.com/spf13/viper"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/internal/output"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:               "in [<project-id>] [<description>]",
	Short:             "Create a new time entry and starts it (will close time entries not closed)",
	Args:              cobra.MaximumNArgs(2),
	ValidArgsFunction: completion.CombineSuggestionsToArgs(suggestWithClientAPI(suggestProjects)),
	Aliases:           []string{"start"},
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error

		tei := dto.TimeEntryImpl{
			WorkspaceID:  viper.GetString(WORKSPACE),
			UserID:       viper.GetString(USER_ID),
			TimeInterval: dto.TimeInterval{},
		}

		if len(args) > 0 {
			tei.ProjectID = args[0]
		}

		if len(args) > 1 {
			tei.Description = args[1]
		}

		if tei, err = fillTimeEntryWithFlags(tei, cmd.Flags()); err != nil {
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

		if err := validateClosingTimeEntry(
			c, tei.WorkspaceID, viper.GetString(USER_ID),
		); err != nil {
			return err
		}

		return manageEntry(
			c,
			tei,
			createTimeEntry(c, viper.GetString(USER_ID), true),
			viper.GetBool(INTERACTIVE),
			getAllowNameForIDsFn(c),
			printTimeEntryImpl(c, cmd, output.TIME_FORMAT_SIMPLE),
			!viper.GetBool(ALLOW_INCOMPLETE),
			true,
			dc,
		)
	}),
}

func init() {
	rootCmd.AddCommand(inCmd)

	addTimeEntryFlags(inCmd)
	addTimeEntryDateFlags(inCmd)
}
