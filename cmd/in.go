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
	"fmt"

	"github.com/spf13/cobra"
)

var cardNumber int
var issueNumber int
var tags []string

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:     "in <project-name-or-id> <description>",
	Short:   "Create a new time entry and starts it",
	Example: `clockify-cli in --issue 13 "timesheet"`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("in called")
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.Flags().IntVarP(&cardNumber, "card", "c", 0, "trello card number being started")
	inCmd.Flags().IntVarP(&issueNumber, "issue", "i", 0, "issue number being started")
	inCmd.Flags().StringSliceVar(&tags, "tag", nil, "add tags to the entry")
	inCmd.Flags().String("when", "", "when the entry should be closed, if not informed will use current time")
}
