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
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/viper"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

var cardNumber int
var issueNumber int
var tags []string
var notBillable bool
var task string

var whenString string
var whenToCloseString string

// inCmd represents the in command
var inCmd = &cobra.Command{
	Use:     "in <project-name-or-id> <description>",
	Short:   "Create a new time entry and starts it",
	Example: `clockify-cli in --issue 13 "timesheet"`,
	Args:    cobra.MaximumNArgs(2),
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {

		var whenDate *time.Time
		var whenToCloseDate *time.Time
		var err error

		workspace := viper.GetString("workspace")
		project := ""
		if len(args) > 0 {
			project = args[0]
		}
		project, err = getProjectID(project, workspace, c)

		if err != nil {
			printError(errors.New("can not end current time entry"))
			return
		}

		if project == "" {
			printError(errors.New("project must be informed"))
			return
		}

		description := getDescription(args, 1)

		tags, err = getTagIDs(tags, workspace, c)
		if err != nil {
			printError(errors.New("can not end current time entry"))
			return
		}

		if whenDate, err = getDateTimeParam("Start", true, whenString, convertToTime); err != nil {
			printError(err)
			return
		}

		if whenToCloseDate, err = getDateTimeParam("End", false, whenToCloseString, convertToTime); err != nil {
			printError(err)
			return
		}

		err = c.Out(api.OutParam{
			Workspace: workspace,
			End:       *whenDate,
		})

		if err != nil {
			printError(errors.New("can not end current time entry"))
			return
		}

		tei, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   workspace,
			Billable:    !notBillable,
			Start:       *whenDate,
			End:         whenToCloseDate,
			ProjectID:   project,
			Description: description,
			TagIDs:      tags,
			TaskID:      task,
		})

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

		var reportFn func(*dto.TimeEntry, io.Writer) error

		reportFn = reports.TimeEntryPrint

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

func getProjectID(value string, workspace string, c *api.Client) (string, error) {
	if value != "" {
		return value, nil
	}

	if !viper.GetBool("interactive") {
		return "", nil
	}

	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace: workspace,
	})

	if err != nil {
		return "", err
	}

	projectsString := make([]string, len(projects))
	for i, u := range projects {
		projectsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
	}

	projectID := ""
	err = survey.AskOne(
		&survey.Select{
			Message: "Choose your project:",
			Options: projectsString,
		},
		&projectID,
		nil,
	)

	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(projectID[0:strings.Index(projectID, " - ")]), nil
}

func getDescription(args []string, i int) string {
	if len(args) > i {
		return args[i]
	}

	if !viper.GetBool("interactive") {
		return ""
	}

	v := ""
	_ = survey.AskOne(
		&survey.Input{
			Message: "Description:",
		},
		&v,
		nil,
	)

	return v
}
func getTagIDs(tagIDs []string, workspace string, c *api.Client) ([]string, error) {
	if len(tagIDs) > 0 {
		return tagIDs, nil
	}

	if !viper.GetBool("interactive") {
		return nil, nil
	}

	tags, err := c.GetTags(api.GetTagsParam{
		Workspace: workspace,
	})

	if err != nil {
		return nil, err
	}

	tagsString := make([]string, len(tags))
	for i, u := range tags {
		tagsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
	}

	err = survey.AskOne(
		&survey.MultiSelect{
			Message: "Choose your tags:",
			Options: tagsString,
		},
		&tagIDs,
		nil,
	)

	if err != nil {
		return nil, nil
	}

	for i, t := range tagIDs {
		tagIDs[i] = strings.TrimSpace(t[0:strings.Index(t, " - ")])
	}

	return tagIDs, nil
}

func init() {
	rootCmd.AddCommand(inCmd)

	addTimeEntryFlags(inCmd)

	inCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each time entry")
	inCmd.Flags().BoolP("json", "j", false, "print as json")
}

func addTimeEntryFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&notBillable, "not-billable", "n", false, "is this time entry not billable")
	cmd.Flags().IntVarP(&cardNumber, "card", "c", 0, "trello card number being started")
	cmd.Flags().IntVar(&issueNumber, "issue", 0, "issue number being started")
	cmd.Flags().StringVar(&task, "task", "", "add a task to the entry")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "add tags to the entry")
	cmd.Flags().StringVar(&whenString, "when", time.Now().Format(fullTimeFormat), "when the entry should be closed, if not informed will use current time")
}
