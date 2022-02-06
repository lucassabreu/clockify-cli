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
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// projectListCmd represents the projectList command
var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects on Clockify and project links",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		asCSV, _ := cmd.Flags().GetBool("csv")
		quiet, _ := cmd.Flags().GetBool("quiet")
		name, _ := cmd.Flags().GetString("name")
		clients, _ := cmd.Flags().GetStringSlice("clients")

		workspace, err := getWorkspaceOrDefault(c)
		if err != nil {
			return err
		}

		if viper.GetBool(ALLOW_NAME_FOR_ID) && len(clients) > 0 {
			clientNames := clients
			clients = make([]string, 0)
			for _, clientName := range clientNames {
				clientId, err := getClientId(clientName, workspace, c)
				if err != nil {
					return err
				}
				clients = append(clients, clientId)
			}
		}

		p := api.GetProjectsParam{
			Workspace:       workspace,
			Name:            name,
			Clients:         clients,
			PaginationParam: api.AllPages(),
		}

		if ok, _ := cmd.Flags().GetBool("not-archived"); ok {
			b := false
			p.Archived = &b
		}
		if ok, _ := cmd.Flags().GetBool("archived"); ok {
			b := true
			p.Archived = &b
		}

		projects, err := c.GetProjects(p)
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

		return reportFn(projects, os.Stdout)
	}),
}

func init() {
	projectCmd.AddCommand(projectListCmd)

	projectListCmd.Flags().StringP("name", "n", "", "will be used to filter the project by name")
	projectListCmd.Flags().StringSliceP("clients", "c", []string{}, "will be used to filter the project by client id/name")
	projectListCmd.Flags().BoolP("not-archived", "", false, "list only active projects")
	projectListCmd.Flags().BoolP("archived", "", false, "list only archived projects")
	projectListCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Project")
	projectListCmd.Flags().BoolP("json", "j", false, "print as JSON")
	projectListCmd.Flags().BoolP("csv", "v", false, "print as CSV")
	projectListCmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
