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
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskAddCmd represents the add command
var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to the specified project",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		f, err := taskReadFlags(cmd)
		if err != nil {
			return err
		}

		project, _ := cmd.Flags().GetString("project")
		if viper.GetBool(ALLOW_NAME_FOR_ID) && project != "" {
			project, err = getProjectByNameOrId(c, f.workspace, project)
			if err != nil {
				return err
			}
		}

		task, err := c.AddTask(api.AddTaskParam{
			Workspace:   f.workspace,
			ProjectID:   project,
			Name:        f.name,
			Estimate:    f.estimate,
			AssigneeIDs: f.assigneeIDs,
			Billable:    f.billable,
		})
		if err != nil {
			return err
		}

		return taskReport(cmd, task)
	}),
}

func init() {
	taskCmd.AddCommand(taskAddCmd)

	addProjectFlags(taskAddCmd)
	taskAddReportFlags(taskAddCmd)
	taskAddPropFlags(taskAddCmd)

	taskAddCmd.MarkFlagRequired("name")
}
