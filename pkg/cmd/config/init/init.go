package init

import (
	"fmt"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func queue(
	tasks ...func() error,
) error {
	for _, t := range tasks {
		if err := t(); err != nil {
			return err
		}
	}

	return nil
}

// NewCmdInit executes and initialization of the config
func NewCmdInit(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Setups the CLI parameters and behavior",
		Long: "Setups the CLI parameters with tokens, default workspace, " +
			"user and behaviors",
		Args: cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
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

			if err := queue(
				func() error { return setWorkspace(c, config, i) },
				func() error { return setUser(c, config, i) },
				updateFlag(
					i, config, cmdutil.CONF_ALLOW_NAME_FOR_ID,
					"Should try to find projects/clients/users/tasks/tags by their names?",
				),
				func() error {
					if !config.IsAllowNameForID() {
						return nil
					}

					return updateFlag(i, config,
						cmdutil.CONF_SEARCH_PROJECTS_WITH_CLIENT_NAME,
						`Should search projects looking into their `+
							`client's name too?`,
					)()
				},
				updateFlag(i, config, cmdutil.CONF_INTERACTIVE,
					`Should use "Interactive Mode" by default?`,
				),
				updateInt(i, config, cmdutil.CONF_INTERACTIVE_PAGE_SIZE,
					"How many items should be shown when asking for "+
						"projects, tasks or tags?"),
				func() error { return setWeekdays(config, i) },
				updateFlag(i, config, cmdutil.CONF_ALLOW_INCOMPLETE,
					`Should allow starting time entries with incomplete data?`,
				),
				updateFlag(i, config, cmdutil.CONF_SHOW_TASKS,
					`Should show task on time entries as a separated column?`,
				),
				updateFlag(i, config, cmdutil.CONF_SHOW_CLIENT,
					`Should show client on time entries as a separated column?`,
				),
				updateFlag(i, config, cmdutil.CONF_SHOW_TOTAL_DURATION,
					`Should show a line with the sum of `+
						`the time entries duration?`,
				),
				updateFlag(i, config, cmdutil.CONF_DESCR_AUTOCOMP,
					`Allow description suggestions using `+
						`recent time entries' descriptions?`,
				),
				func() error {
					if !config.GetBool(cmdutil.CONF_DESCR_AUTOCOMP) {
						config.SetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS, 0)
						return nil
					}
					return updateInt(
						i, config, cmdutil.CONF_DESCR_AUTOCOMP_DAYS,
						`How many days should be used for a time entry to be `+
							`"recent"?`,
					)()
				},
				updateFlag(i, config, cmdutil.CONF_ALLOW_ARCHIVED_TAGS,
					"Should suggest and allow creating time entries "+
						"with archived tags?",
				),
				updateFlag(i, config, cmdutil.CONF_TIME_ENTRY_DEFAULTS,
					"Look for default parameters for time entries per folder?",
					ui.WithConfirmHelp(
						"This will set the default parameters of a time "+
							"entry when using `clockify-cli in` and "+
							"`clockify-cli manual` to the closest "+
							".clockify-defaults.yaml file looking up the "+
							"current folder you were running the commands.\n"+
							"For more information and examples go to "+
							"https://clockify-cli.netlify.app/",
					),
				),
				setLanguage(i, config),
				setTimezone(i, config),
			); err != nil {
				return err
			}

			return config.Save()
		},
	}

	return cmd
}

func setTimezone(i ui.UI, config cmdutil.Config) func() error {
	return func() error {
		tzname, err := i.AskForValidText("What is your preferred timezone:",
			func(s string) error {
				_, err := time.LoadLocation(s)
				return err
			},
			ui.WithHelp("Should be 'Local' to use the systems timezone, UTC "+
				"or valid TZ identifier from the IANA TZ database "+
				"https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"),
			ui.WithDefault(config.TimeZone().String()),
		)
		if err != nil {
			return err
		}

		tz, _ := time.LoadLocation(tzname)

		config.SetTimeZone(tz)

		return nil
	}
}

func setLanguage(i ui.UI, config cmdutil.Config) func() error {
	return func() error {
		suggestLanguages := []string{
			language.English.String(),
			language.German.String(),
			language.Afrikaans.String(),
			language.Chinese.String(),
			language.Portuguese.String(),
		}

		lang, err := i.AskForValidText("What is your preferred language:",
			func(s string) error {
				_, err := language.Parse(s)
				return err
			},
			ui.WithHelp("Accepts any IETF language tag "+
				"https://en.wikipedia.org/wiki/IETF_language_tag"),
			ui.WithSuggestion(func(toComplete string) []string {
				return strhlp.Filter(
					strhlp.IsSimilar(toComplete),
					suggestLanguages,
				)
			}),
			ui.WithDefault(config.Language().String()),
		)
		if err != nil {
			return err
		}

		config.SetLanguage(language.MustParse(lang))
		return nil
	}
}

func setWeekdays(config cmdutil.Config, i ui.UI) (err error) {
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
	return nil
}

func setUser(c api.Client, config cmdutil.Config, i ui.UI) error {
	users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace:       config.GetString(cmdutil.CONF_WORKSPACE),
		PaginationParam: api.AllPages(),
	})

	if err != nil {
		return err
	}

	userID := config.GetString(cmdutil.CONF_USER_ID)
	dUser := ""
	usersString := make([]string, len(users))
	for i := range users {
		usersString[i] = fmt.Sprintf("%s - %s", users[i].ID, users[i].Name)

		if users[i].ID == userID {
			dUser = usersString[i]
		}
	}

	if userID, err = i.AskFromOptions(
		"Choose your user:", usersString, dUser); err != nil {
		return err
	}

	config.SetString(cmdutil.CONF_USER_ID,
		strings.TrimSpace(userID[0:strings.Index(userID, " - ")]))
	return nil
}

func setWorkspace(c api.Client, config cmdutil.Config, i ui.UI) error {
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
	return err
}

func updateInt(ui ui.UI, config cmdutil.Config, param, desc string,
) func() error {
	return func() error {
		value := config.GetInt(param)
		value, err := ui.AskForInt(desc, value)
		if err != nil {
			return err
		}
		config.SetInt(param, value)
		return nil
	}
}

func updateFlag(
	ui ui.UI, config cmdutil.Config, param, description string,
	opts ...ui.ConfirmOption) func() error {
	return func() (err error) {
		b := config.GetBool(param)
		if b, err = ui.Confirm(description, b, opts...); err != nil {
			return
		}
		config.SetBool(param, b)
		return
	}
}
