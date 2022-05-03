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

// taskDeleteCmd represents the close command
var taskDeleteCmd = &cobra.Command{
	Use:     "delete <project> <task>",
	Aliases: []string{"remove", "rm", "del"},
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: completion.CombineSuggestionsToArgs(
		suggestWithClientAPI(suggestProjects)),
	Short: "Deletes a task from a project",
	Long: "Deletes a task from a project\n" +
		"Time entries using this task will revert to not having a task",
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

		p := api.DeleteTaskParam{
			Workspace: workspace,
			ProjectID: project,
			TaskID:    task,
		}

		t, err := c.DeleteTask(p)
		if err != nil {
			return err
		}

		return taskReport(cmd, t)
	}),
}

func init() {
	taskCmd.AddCommand(taskDeleteCmd)

	taskAddReportFlags(taskDeleteCmd)
}
