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
package add

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/project"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdAdd represents the add command
func NewCmdAdd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "create"},
		Short:   "Adds a project to the Clockify workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			asJSON, _ := cmd.Flags().GetBool("json")
			asCSV, _ := cmd.Flags().GetBool("csv")
			quiet, _ := cmd.Flags().GetBool("quiet")
			name, _ := cmd.Flags().GetString("name")
			note, _ := cmd.Flags().GetString("note")
			color, _ := cmd.Flags().GetString("color")
			randomColor, _ := cmd.Flags().GetBool("random-color")
			client, _ := cmd.Flags().GetString("client")
			public, _ := cmd.Flags().GetBool("public")
			billable, _ := cmd.Flags().GetBool("billable")

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() && client != "" {
				cs, err := search.GetClientsByName(c, w, []string{client})
				if err != nil {
					return err
				}
				client = cs[0]
			}

			if randomColor {
				bytes := make([]byte, 3)
				if _, err := rand.Read(bytes); err != nil {
					return err
				}
				color = "#" + hex.EncodeToString(bytes)
			}

			project, err := c.AddProject(api.AddProjectParam{
				Workspace: w,
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

			return report(project, cmd.OutOrStdout(),
				format, asJSON, asCSV, quiet)
		},
	}

	cmd.Flags().StringP("name", "n", "", "name of the new project")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringP("color", "c", "", "color of the new project")
	cmd.Flags().Bool("random-color", false, "use a random color for the project")
	cmd.Flags().StringP("note", "N", "", "note for the new project")
	cmd.Flags().String("client", "", "the id/name of the client the new project will go under")
	cmd.Flags().BoolP("public", "p", false, "make the new project public")
	cmd.Flags().BoolP("billable", "b", false, "make the new project as billable")

	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Project")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")

	return cmd
}

func report(project dto.Project, out io.Writer,
	format string, asJSON, asCSV, quiet bool) error {
	if asJSON {
		return output.ProjectJSONPrint(project, out)
	}

	list := []dto.Project{project}
	if asCSV {
		return output.ProjectsCSVPrint(list, out)
	}

	if format != "" {
		return output.ProjectPrintWithTemplate(format)(list, out)
	}

	if quiet {
		return output.ProjectPrintQuietly(list, out)
	}

	return output.ProjectPrint(list, os.Stdout)
}
