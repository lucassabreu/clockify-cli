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

package user

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/me"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/user"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// NewCmdUser represents the users command
func NewCmdUser(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"user"},
		Short:   "List all users on a Workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			email, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			}
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			users, err := run(c, w, email)
			if err != nil {
				return err
			}

			return report(users, cmd.OutOrStdout(), format, quiet)
		},
	}

	cmd.Flags().StringP("email", "e", "",
		"will be used to filter the workspaces by email")
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each workspace")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	_ = cmd.MarkFlagRequired("workspace")

	cmd.AddCommand(me.NewCmdMe(f))

	return cmd

}

func run(c *api.Client, w, email string) ([]dto.User, error) {
	return c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace: w,
		Email:     email,
	})
}

func report(u []dto.User, out io.Writer, format string, quiet bool) error {
	switch {
	case format != "":
		return output.UserPrintWithTemplate(format)(u, out)
	case quiet:
		return output.UserPrintQuietly(u, out)
	default:
		return output.UserPrint(u, out)
	}
}
