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
			i := f.UI()
			config := f.Config()

			var err error
			token := ""
			if token, err = i.AskForText("User Generated Token:",
				ui.WithDefault(config.GetString(cmdutil.CONF_TOKEN)),
				ui.WithHelp("Can be generated in the following like, "+
					"in the API section: "+
					"https://clockify.me/user/settings#generateApiKeyBtn"),
			); err != nil {
				return err
			}
			config.SetString(cmdutil.CONF_TOKEN, token)

			c, err := f.Client()
			if err != nil {
				return err
			}

			ws, err := c.GetWorkspaces(api.GetWorkspaces{})
			if err != nil {
				return err
			}

			dWorkspace := ""
			wsString := make([]string, len(ws))
			for i := range ws {
				wsString[i] = fmt.Sprintf("%s - %s", ws[i].ID, ws[i].Name)

				if ws[i].ID == config.GetString(cmdutil.CONF_WORKSPACE) {
					dWorkspace = wsString[i]
				}
			}

			w := ""
			if w, err = i.AskFromOptions("Choose default Workspace:",
				wsString, dWorkspace); err != nil {
				return err
			}
			config.SetString(cmdutil.CONF_WORKSPACE,
				strings.TrimSpace(w[0:strings.Index(w, " - ")]))

			users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
				Workspace:       config.GetString(cmdutil.CONF_WORKSPACE),
				PaginationParam: api.AllPages(),
			})

			if err != nil {
				return err
			}

			userId := config.GetString(cmdutil.CONF_USER_ID)
			dUser := ""
			usersString := make([]string, len(users))
			for i := range users {
				usersString[i] = fmt.Sprintf("%s - %s", users[i].ID, users[i].Name)

				if users[i].ID == userId {
					dUser = usersString[i]
				}
			}

			userID := ""
			if userID, err = i.AskFromOptions(
				"Choose your user:", usersString, dUser); err != nil {
				return err
			}
			config.SetString(cmdutil.CONF_USER_ID,
				strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))

			if err := updateFlag(i, config, cmdutil.CONF_ALLOW_NAME_FOR_ID,
				"Should try to find projects/clients/users/tasks/tags by their names?",
			); err != nil {
				return err
			}

			if err := updateFlag(i, config, cmdutil.CONF_INTERACTIVE,
				`Should use "Interactive Mode" by default?`,
			); err != nil {
				return err
			}

			if err = updateInt(i, config, cmdutil.CONF_INTERACTIVE_PAGE_SIZE,
				"How many items should be shown when asking for "+
					"projects, tasks or tags?"); err != nil {
				return err
			}

			workweekDays := config.GetStringSlice(cmdutil.CONF_WORKWEEK_DAYS)
			if workweekDays, err = i.AskManyFromOptions(
				"Which days of the week do you work?",
				cmdutil.GetWeekdays(),
				workweekDays,
				nil,
			); err != nil {
				return err
			}
			config.SetStringSlice(cmdutil.CONF_WORKWEEK_DAYS, workweekDays)

			if err := updateFlag(i, config, cmdutil.CONF_ALLOW_INCOMPLETE,
				`Should allow starting time entries with incomplete data?`,
			); err != nil {
				return err
			}

			if err := updateFlag(i, config, cmdutil.CONF_SHOW_TASKS,
				`Should show task on time entries as a separated column?`,
			); err != nil {
				return err
			}

			if err := updateFlag(i, config, cmdutil.CONF_SHOW_TOTAL_DURATION,
				`Should show a line with the sum of `+
					`the time entries duration?`,
			); err != nil {
				return err
			}

			if err := updateFlag(i, config, cmdutil.CONF_DESCR_AUTOCOMP,
				`Allow description suggestions using `+
					`recent time entries' descriptions?`,
			); err != nil {
				return err
			}

			if config.GetBool(cmdutil.CONF_DESCR_AUTOCOMP) {
				if err := updateInt(
					i, config, cmdutil.CONF_DESCR_AUTOCOMP_DAYS,
					`How many days should be used for a time entry to be `+
						`"recent"?`,
				); err != nil {
					return err
				}
			} else {
				config.SetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS, 0)
			}

			if err := updateFlag(i, config, cmdutil.CONF_ALLOW_ARCHIVED_TAGS,
				"Should suggest and allow creating time entries "+
					"with archived tags?",
			); err != nil {
				return err
			}

			if err := updateFlag(i, config, cmdutil.CONF_TIME_ENTRY_DEFAULTS,
				"Look for default parameters for time entries per folder?",
				ui.WithConfirmHelp(
					"This will set the default parameters of a time entry "+
						"when using `clockify-cli in` and `clockify-cli "+
						"manual` to the closest .clockify-defaults.yaml file "+
						"looking up the current folder you were running the"+
						" commands.\n"+
						"For more information and examples go to "+
						"https://clockify-cli.netlify.app/",
				),
			); err != nil {
				return err
			}

			config.SetString(cmdutil.CONF_TOKEN, token)

			return config.Save()
		},
	}

	return cmd
}

func updateInt(ui ui.UI, config cmdutil.Config, param, desc string) error {
	value := config.GetInt(param)
	value, err := ui.AskForInt(desc, value)
	if err != nil {
		return err
	}
	config.SetInt(param, value)
	return nil
}

func updateFlag(
	ui ui.UI, config cmdutil.Config, param, description string,
	opts ...ui.ConfirmOption) (err error) {
	b := config.GetBool(param)
	if b, err = ui.Confirm(description, b, opts...); err != nil {
		return
	}
	config.SetBool(param, b)
	return
}
