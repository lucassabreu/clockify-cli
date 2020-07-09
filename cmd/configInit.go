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
	"fmt"
	"path"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
)

// configInitCmd represents the configInit command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the config with Tokens, default Workspace and User",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := ""
		err := survey.AskOne(
			&survey.Input{
				Message: "User Generated Token:",
				Default: viper.GetString("token"),
			},
			&token,
			nil,
		)

		if err != nil {
			return err
		}

		viper.Set("token", token)

		c, err := getAPIClient()
		if err != nil {
			return err
		}

		ws, err := c.Workspaces(api.WorkspacesFilter{})

		if err != nil {
			return err
		}

		dWorkspace := ""
		wsString := make([]string, len(ws))
		for i, w := range ws {
			wsString[i] = fmt.Sprintf("%s - %s", w.ID, w.Name)

			if w.ID == viper.GetString("workspace") {
				dWorkspace = wsString[i]
			}
		}

		workspace := ""
		err = survey.AskOne(
			&survey.Select{
				Message: "Choose default Workspace:",
				Options: wsString,
				Default: dWorkspace,
			},
			&workspace,
			nil,
		)

		if err != nil {
			return err
		}

		viper.Set("workspace", strings.TrimSpace(workspace[0:strings.Index(workspace, " - ")]))

		users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
			Workspace: viper.GetString("workspace"),
		})

		if err != nil {
			return err
		}

		dUser := ""
		usersString := make([]string, len(users))
		for i, u := range users {
			usersString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)

			if u.ID == viper.GetString("user.id") {
				dUser = usersString[i]
			}
		}

		userID := ""
		err = survey.AskOne(
			&survey.Select{
				Message: "Choose your user:",
				Options: usersString,
				Default: dUser,
			},
			&userID,
			nil,
		)

		if err != nil {
			return err
		}

		viper.Set("user.id", strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))

		githubToken := ""
		err = survey.AskOne(
			&survey.Input{
				Message: "GitHub token (must have permission to read issues):",
				Default: viper.GetString("github.token"),
			},
			&githubToken,
			nil,
		)

		if err != nil {
			return err
		}

		viper.Set("github.token", githubToken)

		trelloToken := ""
		err = survey.AskOne(
			&survey.Input{
				Message: "Trello token (must have permission to read cards):",
				Default: viper.GetString("trello.token"),
			},
			&trelloToken,
			nil,
		)

		if err != nil {
			return err
		}

		viper.Set("trello.token", trelloToken)

		return saveConfigFile()
	},
}

func saveConfigFile() error {
	filename := viper.ConfigFileUsed()
	if filename == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		filename = path.Join(home, ".clockify-cli.yaml")
	}

	return viper.WriteConfigAs(filename)
}

func init() {
	configCmd.AddCommand(configInitCmd)
}
