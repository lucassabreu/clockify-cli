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

// projectAddCmd represents the add command
var projectAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"new", "create"},
	Short:   "Adds a project to the Clockify workspace",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		asCSV, _ := cmd.Flags().GetBool("csv")
		quiet, _ := cmd.Flags().GetBool("quiet")
		name, _ := cmd.Flags().GetString("name")
		note, _ := cmd.Flags().GetString("note")
		color, _ := cmd.Flags().GetString("color")
		client, _ := cmd.Flags().GetString("client")
		public, _ := cmd.Flags().GetBool("public")
		billable, _ := cmd.Flags().GetBool("billable")

		workspace := viper.GetString(WORKSPACE)

		var err error
		if viper.GetBool(ALLOW_NAME_FOR_ID) && client != "" {
			client, err = getClientByNameOrId(c, workspace, client)
			if err != nil {
				return err
			}
		}

		project, err := c.AddProject(api.AddProjectParam{
			Workspace: workspace,
			Name:      name,
			ClientId:  client,
			Color:     color,
			Billable:  billable,
			Public:    public,
			Note:      note,
		})
		if err != nil {
			return err
		}

		var reportFn func([]dto.Project, io.Writer) error

		reportFn = output.ProjectPrint
		if asJSON {
			reportFn = output.ProjectsJSONPrint
		}

		if asCSV {
			reportFn = output.ProjectsCSVPrint
		}

		if format != "" {
			reportFn = output.ProjectPrintWithTemplate(format)
		}

		if quiet {
			reportFn = output.ProjectPrintQuietly
		}

		return reportFn([]dto.Project{project}, os.Stdout)
	}),
}

func init() {
	projectCmd.AddCommand(projectAddCmd)

	projectAddCmd.Flags().StringP("name", "n", "", "name of the new project")
	projectAddCmd.MarkFlagRequired("name")

	projectAddCmd.Flags().StringP("color", "c", "", "color of the new project")
	projectAddCmd.Flags().StringP("note", "N", "", "note for the new project")
	projectAddCmd.Flags().String("client", "", "the id/name of the client the new project will go under")
	projectAddCmd.Flags().BoolP("public", "p", false, "make the new project public")
	projectAddCmd.Flags().BoolP("billable", "b", false, "make the new project as billable")

	projectAddCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Project")
	projectAddCmd.Flags().BoolP("json", "j", false, "print as JSON")
	projectAddCmd.Flags().BoolP("csv", "v", false, "print as CSV")
}
