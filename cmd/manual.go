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
	"fmt"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// manualCmd represents the manual command
var manualCmd = &cobra.Command{
	Use:               "manual [<project-id>] [<start>] [<end>] [<description>]",
	Short:             "Creates a new completed time entry (does not stop on-going time entries)",
	Args:              cobra.MaximumNArgs(4),
	ValidArgsFunction: completion.CombineSuggestionsToArgs(suggestWithClientAPI(suggestProjects)),
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var whenToCloseDate time.Time
		var err error

		tei := dto.TimeEntryImpl{
			WorkspaceID:  viper.GetString(WORKSPACE),
			TimeInterval: dto.TimeInterval{},
		}

		if len(args) > 0 {
			tei.ProjectID = args[0]
		}

		if len(args) > 1 {
			tei.TimeInterval.Start, err = convertToTime(args[1])
			if err != nil {
				return fmt.Errorf("Fail to convert when to start: %s", err.Error())

			}
		}

		if len(args) > 2 {
			whenToCloseDate, err = convertToTime(args[2])
			if err != nil {
				return fmt.Errorf("Fail to convert when to end: %s", err.Error())
			}
			tei.TimeInterval.End = &whenToCloseDate
		}

		if len(args) > 3 {
			tei.Description = args[3]
		}

		if tei, err = fillTimeEntryWithFlags(tei, cmd.Flags()); err != nil {
			return err
		}

		return manageEntry(
			c,
			tei,
			createTimeEntry(c, viper.GetString(USER_ID), false),
			viper.GetBool(INTERACTIVE),
			viper.GetBool(ALLOW_NAME_FOR_ID),
			printTimeEntryImpl(c, cmd),
			true,
			true,
		)
	}),
}

func init() {
	rootCmd.AddCommand(manualCmd)

	addTimeEntryFlags(manualCmd)
	addTimeEntryDateFlags(manualCmd)
}
