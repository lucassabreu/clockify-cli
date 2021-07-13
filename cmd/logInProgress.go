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
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logInProgressCmd represents the logInProgress command
var logInProgressCmd = &cobra.Command{
	Use:     "in-progress",
	Aliases: []string{"current", "open", "running"},
	Short:   "Show time entry in progress (if any)",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		te, err := getTimeEntryInProgres(c, viper.GetString(WORKSPACE))
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		return formatTimeEntry(te, asJSON, format)
	}),
}

func getTimeEntryInProgres(c *api.Client, workspace string) (*dto.TimeEntry, error) {
	tei, err := c.LogInProgress(api.LogInProgressParam{
		Workspace: workspace,
	})

	if err != nil || tei == nil {
		return nil, err
	}

	tef, err := c.ConvertIntoFullTimeEntry(*tei)
	return &tef, err
}

func formatTimeEntry(te *dto.TimeEntry, asJSON bool, format string) error {
	var reportFn func(*dto.TimeEntry, io.Writer) error

	reportFn = output.TimeEntryPrint

	if asJSON {
		reportFn = output.TimeEntryJSONPrint
	}

	if format != "" {
		reportFn = output.TimeEntryPrintWithTemplate(format)
	}

	return reportFn(te, os.Stdout)
}

func init() {
	logCmd.AddCommand(logInProgressCmd)
	setCommonFormats(logInProgressCmd)

	_ = logInProgressCmd.MarkFlagRequired(WORKSPACE)
}

func setCommonFormats(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied")
	cmd.Flags().BoolP("json", "j", false, "print as json")
}
