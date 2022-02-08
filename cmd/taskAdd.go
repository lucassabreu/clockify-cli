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
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskAddCmd represents the add command
var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to the specified project",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		asCSV, _ := cmd.Flags().GetBool("csv")
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

		p := api.AddTaskParam{
			Workspace: workspace,
			ProjectID: project,
			Name:      name,
		}

		task, err := c.AddTask(p)
		if err != nil {
			return err
		}

		var reportFn func([]dto.Task, io.Writer) error

		reportFn = output.TaskPrint
		if asJSON {
			reportFn = output.TasksJSONPrint
		}

		if asCSV {
			reportFn = output.TasksCSVPrint
		}

		if format != "" {
			reportFn = output.TaskPrintWithTemplate(format)
		}

		return reportFn([]dto.Task{task}, os.Stdout)
	}),
}

func init() {
	taskCmd.AddCommand(taskAddCmd)

	addProjectFlags(taskAddCmd)

	taskAddCmd.Flags().StringP("name", "n", "", "name of the new task")
	taskAddCmd.MarkFlagRequired("name")

	taskAddCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Project")
	taskAddCmd.Flags().BoolP("json", "j", false, "print as JSON")
	taskAddCmd.Flags().BoolP("csv", "v", false, "print as CSV")
}
