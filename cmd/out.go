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
	"errors"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// outCmd represents the out command
var outCmd = &cobra.Command{
	Use:   "out",
	Short: "Stops the last time entry",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		var whenDate time.Time
		var err error
		whenString, _ := cmd.Flags().GetString("when")

		if whenDate, err = convertToTime(whenString); err != nil {
			printError(err)
			return
		}

		workspace := viper.GetString("workspace")
		te, err := getTimeEntryInProgres(c, workspace)
		if te == nil && err == nil {
			err = errors.New("no time entry in progress")
		}

		if err != nil {
			printError(err)
			return
		}

		err = c.Out(api.OutParam{
			Workspace: workspace,
			End:       whenDate,
		})

		if err != nil {
			printError(err)
			return
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			return
		}

		te.TimeInterval.End = &whenDate
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")

		if err = formatTimeEntry(te, asJSON, format); err != nil {
			printError(err)
		}
	}),
}

func init() {
	rootCmd.AddCommand(outCmd)
	setCommonFormats(outCmd)

	outCmd.Flags().String("when", time.Now().Format(fullTimeFormat), "when the entry should be closed, if not informed will use current time")
	outCmd.Flags().BoolP("quiet", "q", false, "print nothing")
}
