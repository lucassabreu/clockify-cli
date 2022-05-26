package init

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/spf13/cobra"
)

func NewCmdInit(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Setups the CLI parameters and behavior",
		Long: "Setups the CLI parameters with tokens, default workspace, " +
			"user and behaviors",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f.Config(), f.Client)
		},
	}

	return cmd
}

func run(config cmdutil.Config, getClient func() (*api.Client, error)) error {
	var err error
	token := ""
	if token, err = ui.AskForText("User Generated Token:",
		ui.WithDefault(config.GetString(cmdutil.CONF_TOKEN)),
		ui.WithHelp("Can be generated here: "+
			"https://clockify.me/user/settings#generateApiKeyBtn")); err != nil {
		return err
	}
	config.SetString(cmdutil.CONF_TOKEN, token)

	c, err := getClient()
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

		if w.ID == config.GetString(cmdutil.CONF_WORKSPACE) {
			dWorkspace = wsString[i]
		}
	}

	workspace := ""
	if workspace, err = ui.AskFromOptions("Choose default Workspace:", wsString, dWorkspace); err != nil {
		return err
	}
	config.SetString(cmdutil.CONF_WORKSPACE,
		strings.TrimSpace(workspace[0:strings.Index(workspace, " - ")]))

	users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace: config.GetString(cmdutil.CONF_WORKSPACE),
	})

	if err != nil {
		return err
	}

	userId := config.GetString(cmdutil.CONF_USER_ID)
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
	config.SetString(cmdutil.CONF_USER_ID,
		strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))

	if err := updateFlag(
		config,
		cmdutil.CONF_ALLOW_NAME_FOR_ID,
		"Should try to find projects/clients/users/tasks/tags by their names?",
	); err != nil {
		return err
	}

	if err := updateFlag(
		config,
		cmdutil.CONF_INTERACTIVE,
		`Should use "Interactive Mode" by default?`,
	); err != nil {
		return err
	}

	workweekDays := config.GetStringSlice(cmdutil.CONF_WORKWEEK_DAYS)
	if workweekDays, err = ui.AskManyFromOptions(
		"Which days of the week do you work?",
		cmdutil.GetWeekdays(),
		workweekDays,
	); err != nil {
		return err
	}
	config.SetStringSlice(cmdutil.CONF_WORKWEEK_DAYS, workweekDays)

	if err := updateFlag(
		config,
		cmdutil.CONF_ALLOW_INCOMPLETE,
		`Should allow starting time entries with incomplete data?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		config,
		cmdutil.CONF_SHOW_TASKS,
		`Should show task on time entries as a separated column?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		config,
		cmdutil.CONF_SHOW_TOTAL_DURATION,
		`Should show a line with the sum of the time entries duration?`,
	); err != nil {
		return err
	}

	if err := updateFlag(
		config,
		cmdutil.CONF_DESCR_AUTOCOMP,
		`Allow description suggestions using recent time entries' descriptions?`,
	); err != nil {
		return err
	}

	daysToConsider := config.GetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS)
	if config.GetBool(cmdutil.CONF_DESCR_AUTOCOMP) {
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

	config.SetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS, daysToConsider)

	return config.Save()
}

func updateFlag(config cmdutil.Config,
	param string, description string) (err error) {
	b := config.GetBool(param)
	if b, err = ui.Confirm(description, b); err != nil {
		return
	}
	config.SetBool(param, b)
	return
}
