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

// usersCmd represents the workspaceUsers command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users on a Workspace",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		email, _ := cmd.Flags().GetString("email")
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")

		users, err := getUsers(c, email)

		if err != nil {
			return err
		}

		var reportFn func([]dto.User, io.Writer) error

		switch true {
		case format != "":
			reportFn = output.UserPrintWithTemplate(format)
		case quiet:
			reportFn = output.UserPrintQuietly
		default:
			reportFn = output.UserPrint
		}

		return reportFn(users, os.Stdout)
	}),
}

func getUsers(c *api.Client, email string) ([]dto.User, error) {
	return c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace: viper.GetString(WORKSPACE),
		Email:     email,
	})
}

func init() {
	workspacesCmd.AddCommand(usersCmd)

	usersCmd.Flags().StringP("email", "e", "", "will be used to filter the workspaces by email")
	usersCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each workspace")
	usersCmd.Flags().BoolP("quiet", "q", false, "only display ids")

	_ = usersCmd.MarkFlagRequired(WORKSPACE)
}
