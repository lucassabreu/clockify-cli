/*
Copyright © 2021 Lucas dos Santos Abreu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:       "show [" + ALIAS_CURRENT + "|" + ALIAS_LAST + "|<time-entry-id>|^n]",
	ValidArgs: []string{ALIAS_CURRENT, ALIAS_LAST},
	Args:      cobra.MaximumNArgs(1),
	Short:     "Show detailed information about one time entry.",
	Long:      "Show detailed information about one time entry.\nShows current one by default",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		id := ALIAS_CURRENT
		if len(args) > 0 {
			id = args[0]
		}

		tei, err := getTimeEntry(
			id,
			viper.GetString(WORKSPACE),
			userID,
			c,
		)

		if err != nil {
			return err
		}

		return printTimeEntryImpl(tei, c, cmd, fullTimeFormat)
	}),
}

func init() {
	rootCmd.AddCommand(showCmd)

	addPrintTimeEntriesFlags(showCmd)
	_ = showCmd.MarkFlagRequired(WORKSPACE)
	_ = showCmd.MarkFlagRequired(USER_ID)
}
