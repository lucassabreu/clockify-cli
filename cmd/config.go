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
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/ui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var configValidArgs = completion.ValigsArgsMap{
	"token":              `clockify's token`,
	"workspace":          "workspace to be used",
	"user.id":            "user id from the token",
	"allow-project-name": "should allow use of project when id is asked",
	"no-closing":         "should not close any active time entry",
	"interactive":        "show interactive mode",
}

const FORMAT_YAML = "yaml"
const FORMAT_JSON = "json"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:       "config " + configValidArgs.IntoUse() + " [value]",
	Short:     "Manages configuration file parameters",
	Args:      cobra.MaximumNArgs(2),
	ValidArgs: configValidArgs.IntoValidArgs(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if b, _ := cmd.Flags().GetBool("init"); b {
			return configInit(cmd, args)
		}

		if len(args) < 2 {
			return configShow(cmd, args)
		}

		return configSet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().Bool("init", false, "initialize the config with tokens, default workspace, user and behaviors")

	configCmd.Flags().StringP("format", "f", FORMAT_YAML, "output format (when not setting or initializing)")
	_ = completion.AddFixedSuggestionsToFlag(configCmd, "format", completion.ValigsArgsSlide{FORMAT_YAML, FORMAT_JSON})
}

func configShow(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")

	var b []byte

	var v interface{}
	if len(args) == 0 {
		v = viper.AllSettings()
	} else {
		v = viper.Get(args[0])
	}

	format = strings.ToLower(format)
	switch format {
	case FORMAT_JSON:
		b, _ = json.Marshal(v)

	case FORMAT_YAML:
		b, _ = yaml.Marshal(v)
	default:
		return errors.New("invalid format")
	}

	fmt.Println(string(b))
	return nil
}

func configSaveFile() error {
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

func configInit(cmd *cobra.Command, args []string) error {
	var err error
	token := ""
	if token, err = ui.AskForText("User Generated Token:", viper.GetString("token")); err != nil {
		return err
	}
	viper.Set("token", token)

	c, err := getAPIClient()
	if err != nil {
		return err
	}

	ws, err := c.GetWorkspaces(api.GetWorkspaces{})
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
	if workspace, err = ui.AskFromOptions("Choose default Workspace:", wsString, dWorkspace); err != nil {
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
	if userID, err = ui.AskFromOptions("Choose your user:", usersString, dUser); err != nil {
		return err
	}
	viper.Set("user.id", strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))

	allowProjectName := viper.GetBool("allow-project-name")
	if allowProjectName, err = ui.Confirm(
		"Should try to find project by its name?",
		allowProjectName,
	); err != nil {
		return err
	}
	viper.Set("allow-project-name", allowProjectName)

	autoClose := !viper.GetBool("no-closing")
	if autoClose, err = ui.Confirm(
		`Should auto-close previous/current time entry before opening a new one?`,
		autoClose,
	); err != nil {
		return err
	}
	viper.Set("no-closing", !autoClose)

	interactive := viper.GetBool("interactive")
	if interactive, err = ui.Confirm(
		`Should use "Interactive Mode" by default?`,
		interactive,
	); err != nil {
		return err
	}
	viper.Set("interactive", interactive)

	return configSaveFile()
}

func configSet(cmd *cobra.Command, args []string) error {
	viper.Set(args[0], args[1])
	return configSaveFile()
}
