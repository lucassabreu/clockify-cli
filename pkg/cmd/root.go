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
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/completion"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/tag"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/me"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/version"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/workspace"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the base command when called without any subcommands
func NewCmdRoot(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "clockify-cli",
		Short:         "Allow to integrate with Clockify through terminal",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().StringP("token", "t", "",
		"clockify's token\nCan be generated here: "+
			"https://clockify.me/user/settings#generateApiKeyBtn")

	cmd.PersistentFlags().StringP("workspace", "w", "", "workspace to be used")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "workspace",
		cmdcomplutil.NewWorspaceAutoComplete(f))

	cmd.PersistentFlags().StringP("user-id", "u", "", "user id from the token")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "user-id",
		cmdcomplutil.NewUserAutoComplete(f))

	cmd.PersistentFlags().Bool("debug", false, "show debug log")

	cmd.PersistentFlags().BoolP("interactive", "i", false,
		"will prompt you to confirm/complement commands input before "+
			"executing the action ")

	cmd.PersistentFlags().BoolP("allow-name-for-id", "", false,
		"allow use of project/client/tag's name when id is asked")

	_ = cmd.MarkFlagRequired("token")

	cmd.AddCommand(version.NewCmdVersion(f))

	cmd.AddCommand(config.NewCmdConfig(f))

	cmd.AddCommand(workspace.NewCmdWorkspace(f))

	cmd.AddCommand(user.NewCmdUser(f))
	cmd.AddCommand(me.NewCmdMe(f))

	cmd.AddCommand(client.NewCmdClient(f))
	cmd.AddCommand(project.NewCmdProject(f))
	cmd.AddCommand(task.NewCmdTask(f))

	cmd.AddCommand(tag.NewCmdTag(f))

	cmd.AddCommand(timeentry.NewCmdTimeEntry(f)...)

	cmd.AddCommand(completion.NewCmdCompletion())

	cmd.AddCommand(gendocsCmd)

	return cmd
}