/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"io"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskCmd represents the client command
var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"tasks"},
	Short:   "List/add tasks of/to a project",
}

func init() {
	rootCmd.AddCommand(taskCmd)
}

func taskReport(cmd *cobra.Command, tasks ...dto.Task) error {
	var reportFn func([]dto.Task, io.Writer) error

	reportFn = output.TaskPrint
	if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
		reportFn = output.TasksJSONPrint
	}

	if asCSV, _ := cmd.Flags().GetBool("csv"); asCSV {
		reportFn = output.TasksCSVPrint
	}

	if format, _ := cmd.Flags().GetString("format"); format != "" {
		reportFn = output.TaskPrintWithTemplate(format)
	}

	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		reportFn = output.TaskPrintQuietly
	}

	return reportFn(tasks, cmd.OutOrStdout())
}

func taskAddReportFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Client")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")
}

func taskAddPropFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "new name of the task")
	cmd.Flags().Int32P("estimate", "E", 0, "estimation on hours")
	cmd.Flags().Bool("billable", false, "sets the task as billable")
	cmd.Flags().Bool("not-billable", false, "sets the task as not billable")

	cmd.Flags().StringSliceP("assignee", "A", []string{},
		"list of users that are assigned to this task")
	_ = completion.AddSuggestionsToFlag(cmd, "assignee",
		suggestWithClientAPI(suggestUsers))

	cmd.Flags().Bool("no-assignee", false,
		"cleans the assignee list")
}

func taskReadFlags(cmd *cobra.Command) (p struct {
	workspace   string
	name        string
	estimate    *time.Duration
	assigneeIDs *[]string
	billable    *bool
}, err error) {
	p.workspace = viper.GetString(WORKSPACE)
	p.name, _ = cmd.Flags().GetString("name")
	if cmd.Flags().Changed("estimate") {
		e, _ := cmd.Flags().GetInt32("estimate")
		d := time.Duration(e) * time.Hour
		p.estimate = &d
	}

	if cmd.Flags().Changed("assignee") {
		assignees, _ := cmd.Flags().GetStringSlice("assignee")
		p.assigneeIDs = &assignees
	}

	if cmd.Flags().Changed("no-assignee") {
		a := []string{}
		p.assigneeIDs = &a
	}

	switch {
	case cmd.Flags().Changed("billable"):
		b := true
		p.billable = &b
	case cmd.Flags().Changed("not-billable"):
		b := false
		p.billable = &b
	}

	return
}
