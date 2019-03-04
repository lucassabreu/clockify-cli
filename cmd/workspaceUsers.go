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
	"github.com/lucassabreu/clockify-cli/reports"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// workspaceUsersCmd represents the workspaceUsers command
var workspaceUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users on a Workspace",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		workspace, _ := cmd.Flags().GetString("workspace")
		email, _ := cmd.Flags().GetString("email")
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")

		users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
			Workspace: workspace,
			Email:     email,
		})

		if err != nil {
			printError(err)
			return
		}

		var reportFn func([]dto.User, io.Writer) error

		reportFn = reports.UserPrint
		if format != "" {
			reportFn = reports.UserPrintWithTemplate(format)
		}

		if quiet {
			reportFn = reports.UserPrintQuietly
		}

		if err = reportFn(users, os.Stdout); err != nil {
			printError(err)
		}

	}),
}

func init() {
	workspacesCmd.AddCommand(workspaceUsersCmd)

	workspaceUsersCmd.Flags().StringP("email", "e", "", "will be used to filter the workspaces by email")
	workspaceUsersCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each workspace")
	workspaceUsersCmd.Flags().BoolP("quiet", "q", false, "only display ids")

	workspaceUsersCmd.MarkFlagRequired("workspace")
}
