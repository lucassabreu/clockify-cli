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
	"io"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/client"
	"github.com/spf13/cobra"
)

// clientAddCmd represents the add command
func NewCmdAdd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "create"},
		Short:   "Adds a client to the Clockify workspace",
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

			client, err := c.AddClient(api.AddClientParam{
				Workspace: w,
				Name:      name,
			})
			if err != nil {
				return err
			}

			return report(client, cmd.OutOrStdout(),
				format, asCSV, asJSON, quiet)
		},
	}

	cmd.Flags().StringP("name", "n", "", "the name of the new client")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each Client")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}

func report(c dto.Client, out io.Writer,
	format string, asCSV, asJSON, quiet bool) error {

	if asJSON {
		return output.ClientJSONPrint(c, out)
	}

	cs := []dto.Client{c}
	if asCSV {
		return output.ClientsCSVPrint(cs, out)
	}

	if format != "" {
		return output.ClientPrintWithTemplate(format)(cs, out)
	}

	if quiet {
		return output.ClientPrintQuietly(cs, out)
	}

	return output.ClientPrint(cs, out)
}
