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
package list

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/client"
	"github.com/spf13/cobra"
)

// NewCmdList represents the list command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List clients on Clockify",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			asJSON, _ := cmd.Flags().GetBool("json")
			asCSV, _ := cmd.Flags().GetBool("csv")
			quiet, _ := cmd.Flags().GetBool("quiet")
			name, _ := cmd.Flags().GetString("name")

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var archived *bool
			if ok, _ := cmd.Flags().GetBool("not-archived"); ok {
				b := false
				archived = &b
			}
			if ok, _ := cmd.Flags().GetBool("archived"); ok {
				b := true
				archived = &b
			}
			clients, err := run(c, w, name, archived)

			if err != nil {
				return err
			}

			return report(clients, cmd.OutOrStdout(),
				format, asCSV, asJSON, quiet)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().BoolP("not-archived", "", false, "list only active projects")
	cmd.Flags().BoolP("archived", "", false, "list only archived projects")
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each Client")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}

func run(c *api.Client, w, name string, archived *bool) ([]dto.Client, error) {

	return c.GetClients(api.GetClientsParam{
		Workspace:       w,
		Name:            name,
		Archived:        archived,
		PaginationParam: api.AllPages(),
	})

}

func report(clients []dto.Client, out io.Writer,
	format string, asCSV, asJSON, quiet bool) error {

	if asJSON {
		return output.ClientsJSONPrint(clients, out)
	}

	if asCSV {
		return output.ClientsCSVPrint(clients, out)
	}

	if format != "" {
		return output.ClientPrintWithTemplate(format)(clients, out)
	}

	if quiet {
		return output.ClientPrintQuietly(clients, out)
	}

	return output.ClientPrint(clients, out)
}
