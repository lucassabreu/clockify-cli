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

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/output"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// workspacesCmd represents the workspaces command
var workspacesCmd = &cobra.Command{
	Use:     WORKSPACE,
	Aliases: []string{"workspaces"},
	Short:   "List user's workspaces",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		name, _ := cmd.Flags().GetString("name")
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")

		w, err := getWorkspaces(c, name)
		if err != nil {
			return err
		}

		var reportFn func([]dto.Workspace, io.Writer) error

		reportFn = output.WorkspacePrint(viper.GetString(WORKSPACE))
		if format != "" {
			reportFn = output.WorkspacePrintWithTemplate(format)
		}

		if quiet {
			reportFn = output.WorkspacePrintQuietly
		}

		return reportFn(w, os.Stdout)
	}),
}

func getWorkspaces(c *api.Client, name string) ([]dto.Workspace, error) {
	return c.GetWorkspaces(api.GetWorkspaces{Name: name})
}

func init() {
	rootCmd.AddCommand(workspacesCmd)

	workspacesCmd.Flags().StringP("name", "n", "", "will be used to filter the workspaces by name")
	workspacesCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each workspace")
	workspacesCmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
