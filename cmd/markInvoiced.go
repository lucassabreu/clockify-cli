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
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// markInvoicedCmd represents the mark invoiced command
var markInvoicedCmd = &cobra.Command{
	Use:       "mark-invoiced [current|last|<time-entry-id>]...",
	Short:     "Marks times entries as invoiced",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"last", "current"},
	RunE:      changeInvoiced(true),
}

// markNotInvoiced represents the mark invoiced command
var markNotInvoiced = &cobra.Command{
	Use:       "mark-not-invoiced [current|last|<time-entry-id>]...",
	Short:     "Mark times entries as not invoiced",
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"last", "current"},
	RunE:      changeInvoiced(false),
}

func changeInvoiced(invoiced bool) func(cmd *cobra.Command, args []string) error {
	return withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		workspace := viper.GetString(WORKSPACE)
		userID := viper.GetString(USER_ID)

		args = strhlp.Unique(args)
		tes := make([]dto.TimeEntry, len(args))
		for i, id := range args {
			if id == "current" || id == "last" {
				tei, err := getTimeEntry(id, workspace, userID, false, c)
				if err != nil {
					return err
				}
				id = tei.ID
				args[i] = id
			}

			te, err := c.GetHydratedTimeEntry(api.GetTimeEntryParam{
				Workspace:   workspace,
				TimeEntryID: id,
			})
			if err != nil {
				return err
			}

			tes[i] = *te
		}

		if err := c.ChangeInvoiced(api.ChangeInvoicedParam{
			Workspace:    workspace,
			TimeEntryIDs: args,
			Invoiced:     invoiced,
		}); err != nil {
			return err
		}

		return printTimeEntries(tes, cmd)
	})
}

func init() {
	rootCmd.AddCommand(markNotInvoiced)
	rootCmd.AddCommand(markInvoicedCmd)
}
