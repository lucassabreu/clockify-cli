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
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/lucassabreu/clockify-cli/ui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

const (
	WORKWEEK_DAYS       = "workweek-days"
	INTERACTIVE         = "interactive"
	ALLOW_NAME_FOR_ID   = "allow-name-for-id"
	USER_ID             = "user.id"
	WORKSPACE           = "workspace"
	TOKEN               = "token"
	ALLOW_INCOMPLETE    = "allow-incomplete"
	SHOW_TASKS          = "show-task"
	DESCR_AUTOCOMP      = "description-autocomplete"
	DESCR_AUTOCOMP_DAYS = "description-autocomplete-days"
	SHOW_TOTAL_DURATION = "show-total-duration"
)

var configValidArgs = completion.ValigsArgsMap{
	TOKEN:               `clockify's token`,
	WORKSPACE:           "workspace to be used",
	USER_ID:             "user id from the token",
	ALLOW_NAME_FOR_ID:   "allow to input the name of the entity instead of its ID (projects and tags)",
	INTERACTIVE:         "show interactive mode",
	WORKWEEK_DAYS:       "days of the week were your expected to work (use comma to set multiple)",
	ALLOW_INCOMPLETE:    "should allow starting time entries with missing required values",
	SHOW_TASKS:          "should show an extra column with the task description",
	DESCR_AUTOCOMP:      "autocomplete description looking at recent time entries",
	DESCR_AUTOCOMP_DAYS: "how many days should be considered for the description autocomplete",
	SHOW_TOTAL_DURATION: "adds a totals line on time entry reports with the sum of the time entries duration",
}

var weekdays []string

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

	weekdays = []string{
		time.Sunday:    strings.ToLower(time.Sunday.String()),
		time.Monday:    strings.ToLower(time.Monday.String()),
		time.Tuesday:   strings.ToLower(time.Tuesday.String()),
		time.Wednesday: strings.ToLower(time.Wednesday.String()),
		time.Thursday:  strings.ToLower(time.Thursday.String()),
		time.Friday:    strings.ToLower(time.Friday.String()),
		time.Saturday:  strings.ToLower(time.Saturday.String()),
	}
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

func configInit(_ *cobra.Command, _ []string) error {
	var err error
	token := ""
	if token, err = ui.AskForText("User Generated Token:", viper.GetString(TOKEN)); err != nil {
		return err
	}
	viper.Set(TOKEN, token)

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

		if w.ID == viper.GetString(WORKSPACE) {
			dWorkspace = wsString[i]
		}
	}

	workspace := ""
	if workspace, err = ui.AskFromOptions("Choose default Workspace:", wsString, dWorkspace); err != nil {
		return err
	}
	viper.Set(WORKSPACE, strings.TrimSpace(workspace[0:strings.Index(workspace, " - ")]))

	users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace: viper.GetString(WORKSPACE),
	})

	if err != nil {
		return err
	}

	userId := viper.GetString(USER_ID)
	dUser := ""
	usersString := make([]string, len(users))
	for i, u := range users {
		usersString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)

		if u.ID == userId {
			dUser = usersString[i]
		}
	}

	userID := ""
	if userID, err = ui.AskFromOptions("Choose your user:", usersString, dUser); err != nil {
		return err
	}
	viper.Set(USER_ID, strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))

	if err := updateFlag(
		ALLOW_NAME_FOR_ID,
		"Should try to find projects/tasks/tags by their names?",
	); err != nil {
		return err
	}

	if err := updateFlag(
		INTERACTIVE,
		`Should use "Interactive Mode" by default?`,
	); err != nil {
		return err
	}

	workweekDays := viper.GetStringSlice(WORKWEEK_DAYS)
	if workweekDays, err = ui.AskManyFromOptions(
		"Which days of the week do you work?",
		weekdays,
		workweekDays,
	); err != nil {
		return err
	}
	viper.Set(WORKWEEK_DAYS, workweekDays)

	if err := updateFlag(
		ALLOW_INCOMPLETE,
		`Should allow starting time entries with incomplete data?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		SHOW_TASKS,
		`Should show task on time entries as a separated column?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		SHOW_TOTAL_DURATION,
		`Should show a line with the sum of the time entries duration?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		DESCR_AUTOCOMP,
		`Allow description suggestions using recent time entries' descriptions?`,
	); err != nil {
		return err
	}

	daysToConsider := viper.GetInt(DESCR_AUTOCOMP_DAYS)
	if viper.GetBool(DESCR_AUTOCOMP) {
		daysToConsider, err = ui.AskForInt(
			`How many days should be used for a time entry to be "recent"?`,
			daysToConsider,
		)
		if err != nil {
			return err
		}
	} else {
		daysToConsider = 0
	}

	viper.Set(DESCR_AUTOCOMP_DAYS, daysToConsider)

	return configSaveFile()
}

func updateFlag(config string, description string) (err error) {
	b := viper.GetBool(config)
	if b, err = ui.Confirm(description, b); err != nil {
		return
	}
	viper.Set(config, b)
	return
}

func configSet(_ *cobra.Command, args []string) error {
	switch args[0] {
	case WORKWEEK_DAYS:
		ws := strings.Split(strings.ToLower(args[1]), ",")
		ws = strhlp.Filter(
			func(s string) bool { return strhlp.Search(s, weekdays) != -1 },
			ws,
		)
		viper.Set(args[0], ws)
	default:
		viper.Set(args[0], args[1])
	}
	return configSaveFile()
}
