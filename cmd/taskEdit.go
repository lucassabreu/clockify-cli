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
	"errors"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskEditCmd represents the close command
var taskEditCmd = &cobra.Command{
	Use:     "edit <project> <task>",
	Aliases: []string{"update"},
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: completion.CombineSuggestionsToArgs(
		suggestWithClientAPI(suggestProjects)),
	Short: "Edit a task from a project",
	RunE: withClockifyClient(func(
		cmd *cobra.Command, args []string, c *api.Client) error {
		if len(args) != 2 {
			return errors.New("two arguments are required (project and task)")
		}

		project := args[0]
		task := args[1]

		workspace := viper.GetString(WORKSPACE)
		var err error
		if viper.GetBool(ALLOW_NAME_FOR_ID) {
			if project != "" {
				project, err = getProjectByNameOrId(c, workspace, project)
				if err != nil {
					return err
				}
			}

			if task != "" {
				task, err = getTaskByNameOrId(c, workspace, project, task)
				if err != nil {
					return err
				}
			}
		}

		f, err := taskReadFlags(cmd)
		if err != nil {
			return err
		}

		p := api.UpdateTaskParam{
			Workspace:   f.workspace,
			ProjectID:   project,
			TaskID:      task,
			Name:        f.name,
			Estimate:    f.estimate,
			AssigneeIDs: f.assigneeIDs,
			Billable:    f.billable,
		}

		if !cmd.Flags().Changed("name") {
			t, err := c.GetTask(api.GetTaskParam{
				Workspace: workspace,
				ProjectID: project,
				TaskID:    task,
			})
			if err != nil {
				return err
			}

			p.Name = t.Name
		}

		switch {
		case cmd.Flags().Changed("active"):
			p.Status = api.TaskStatusActive
		case cmd.Flags().Changed("done"):
			p.Status = api.TaskStatusDone
		default:
			p.Status = api.TaskStatusDefault
		}

		t, err := c.UpdateTask(p)
		if err != nil {
			return err
		}

		return taskReport(cmd, t)
	}),
}

func init() {
	taskCmd.AddCommand(taskEditCmd)

	taskAddPropFlags(taskEditCmd)
	taskEditCmd.Flags().Bool("done", false, "sets the task as done")
	taskEditCmd.Flags().Bool("active", false, "sets the task as active")

	taskAddReportFlags(taskEditCmd)
}
