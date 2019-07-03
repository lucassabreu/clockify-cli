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
	"os"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit <time-entry-id>",
	Aliases: []string{"update"},
	Args:    cobra.ExactArgs(1),
	Short:   `Edit a time entry, use id "current" to apply to time entry in progress`,
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		var err error
		param := api.UpdateTimeEntryParam{
			Workspace:   viper.GetString("workspace"),
			TimeEntryID: args[0],
		}

		param, err = getForEdit(param, c)
		if err != nil {
			printError(err)
			return
		}

		if cmd.Flags().Changed("not-billable") {
			b, _ := cmd.Flags().GetBool("not-billable")
			param.Billable = !b
		}

		if cmd.Flags().Changed("project") {
			param.ProjectID, _ = cmd.Flags().GetString("project")
		}

		if cmd.Flags().Changed("description") {
			param.Description, _ = cmd.Flags().GetString("description")
		}

		if cmd.Flags().Changed("task") {
			param.TaskID, _ = cmd.Flags().GetString("task")
		}

		if cmd.Flags().Changed("tag") {
			param.TagIDs, _ = cmd.Flags().GetStringSlice("tag")
		}

		if !cmd.Flags().Changed("when") {
			whenString, _ = cmd.Flags().GetString("when")
			if param.Start, err = convertToTime(whenString); err != nil {
				printError(err)
				return
			}
		}

		if cmd.Flags().Changed("end-at") {
			whenString, _ = cmd.Flags().GetString("end-at")
			var v time.Time
			if v, err = convertToTime(whenString); err != nil {
				printError(err)
				return
			}
			param.End = &v
		}

		if viper.GetBool("interactive") {
			if param, err = confirmValuesForEdit(param, c); err != nil {
				printError(err)
				return
			}
		}

		tei, err := c.UpdateTimeEntry(param)

		if err != nil {
			printError(err)
			return
		}

		te, err := c.ConvertIntoFullTimeEntry(tei)
		if err != nil {
			printError(err)
			return
		}

		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")

		reportFn := reports.TimeEntryPrint

		if asJSON {
			reportFn = reports.TimeEntryJSONPrint
		}

		if format != "" {
			reportFn = reports.TimeEntryPrintWithTemplate(format)
		}

		if err = reportFn(&te, os.Stdout); err != nil {
			printError(err)
		}
	}),
}

func confirmValuesForEdit(param api.UpdateTimeEntryParam, c *api.Client) (api.UpdateTimeEntryParam, error) {
	var err error
	if param.ProjectID, err = getProjectID(param.ProjectID, param.Workspace, c); err != nil {
		return param, err
	}

	if param.ProjectID == "" {
		return param, errors.New("project must be informed")
	}

	_ = survey.AskOne(
		&survey.Input{
			Message: "Description:",
			Default: param.Description,
		},
		&param.Description,
		nil,
	)

	param.TagIDs, err = getTagIDs(param.TagIDs, param.Workspace, c)
	if err != nil {
		return param, errors.New("can not end current time entry")
	}

	var t *time.Time
	if t, err = getDateTimeParam("Start", true, param.Start.Format(fullTimeFormat), convertToTime); err != nil {
		return param, err
	}
	param.Start = *t

	when := ""
	if param.End != nil {
		when = param.End.Format(fullTimeFormat)
	}

	if param.End, err = getDateTimeParam("End", false, when, convertToTime); err != nil {
		return param, err
	}

	return param, nil
}

func getForEdit(param api.UpdateTimeEntryParam, c *api.Client) (api.UpdateTimeEntryParam, error) {
	if param.TimeEntryID == "current" {
		te, err := c.LogInProgress(api.LogInProgressParam{
			Workspace: param.Workspace,
		})

		if err != nil {
			return param, err
		}

		if te == nil {
			return param, errors.New("there is no time entry in progress")
		}

		param.TimeEntryID = te.ID

		if !viper.GetBool("interactive") {
			return param, nil
		}

		param.ProjectID = te.ProjectID
		param.Description = te.Description
		param.TaskID = te.TaskID
		param.TagIDs = te.TagIDs
		param.Billable = te.Billable
		param.Start = te.TimeInterval.Start
		param.End = te.TimeInterval.End
		return param, nil
	}

	tec, err := getTimeEntry(
		param.TimeEntryID,
		param.Workspace,
		viper.GetString("user.id"),
		c,
	)

	if err != nil {
		return param, err
	}

	if !viper.GetBool("interactive") {
		return param, nil
	}

	param.ProjectID = tec.ProjectID
	param.Description = tec.Description
	param.TaskID = tec.TaskID
	param.TagIDs = tec.TagIDs
	param.Billable = tec.Billable
	param.Start = tec.TimeInterval.Start
	param.End = tec.TimeInterval.End
	return param, nil
}

func init() {
	rootCmd.AddCommand(editCmd)

	addTimeEntryFlags(editCmd)

	editCmd.Flags().StringP("project", "p", "", "change the project")
	editCmd.Flags().String("description", "", "change the description")
	editCmd.Flags().String("end-at", "", "when the entry should end (if not set \"\" will be used)")

	editCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	editCmd.Flags().BoolP("json", "j", false, "print as json")
}
