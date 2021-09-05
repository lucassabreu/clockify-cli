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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logInProgressCmd represents the logInProgress command
var logInProgressCmd = &cobra.Command{
	Use:     "in-progress",
	Aliases: []string{"current", "open", "running"},
	Short:   "Show time entry in progress (if any)",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		te, err := c.GetHydratedTimeEntryInProgress(api.GetTimeEntryInProgressParam{
			Workspace: viper.GetString(WORKSPACE),
			UserID:    viper.GetString(USER_ID),
		})

		if err != nil {
			return err
		}

		return formatTimeEntry(te, cmd)
	}),
}

func init() {
	logCmd.AddCommand(logInProgressCmd)
	addPrintTimeEntriesFlags(logInProgressCmd)

	_ = logInProgressCmd.MarkFlagRequired(WORKSPACE)
}
