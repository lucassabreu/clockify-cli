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

// taskListCmd represents the list command
var taskListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List tasks of a Clockify project",
	Aliases: []string{"ls"},
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		active, _ := cmd.Flags().GetBool("active")
		name, _ := cmd.Flags().GetString("name")
		project, _ := cmd.Flags().GetString("project")

		workspace := viper.GetString(WORKSPACE)

		var err error
		if viper.GetBool(ALLOW_NAME_FOR_ID) && project != "" {
			project, err = getProjectByNameOrId(c, workspace, project)
			if err != nil {
				return err
			}
		}

		p := api.GetTasksParam{
			Workspace:       workspace,
			Name:            name,
			ProjectID:       project,
			Active:          active,
			PaginationParam: api.AllPages(),
		}

		tasks, err := c.GetTasks(p)
		if err != nil {
			return err
		}

		return taskReport(cmd, tasks...)
	}),
}

func init() {
	taskCmd.AddCommand(taskListCmd)

	addProjectFlags(taskListCmd)
	taskAddReportFlags(taskListCmd)

	taskListCmd.Flags().StringP("name", "n", "", "will be used to filter the tag by name")
	taskListCmd.Flags().BoolP("active", "a", false, "display only active tasks")
}
