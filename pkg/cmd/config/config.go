package config

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/get"
	initialize "github.com/lucassabreu/clockify-cli/pkg/cmd/config/init"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"

	"github.com/spf13/cobra"
)

var validParameters = cmdcompl.ValidArgsMap{
	cmdutil.CONF_TOKEN:     "clockify's token",
	cmdutil.CONF_WORKSPACE: "workspace to be used",
	cmdutil.CONF_USER_ID:   "user id from the token",
	cmdutil.CONF_ALLOW_NAME_FOR_ID: "allow to input the name of the entity " +
		"instead of its ID (projects, clients, tasks, users and tags)",
	cmdutil.CONF_INTERACTIVE: "show interactive mode",
	cmdutil.CONF_WORKWEEK_DAYS: "days of the week were your expected to " +
		"work (use comma to set multiple)",
	cmdutil.CONF_ALLOW_INCOMPLETE: "should allow starting time entries with " +
		"missing required values",
	cmdutil.CONF_SHOW_TASKS: "should show an extra column with the task " +
		"description",
	cmdutil.CONF_DESCR_AUTOCOMP: "autocomplete description looking at " +
		"recent time entries",
	cmdutil.CONF_DESCR_AUTOCOMP_DAYS: "how many days should be considered " +
		"for the description autocomplete",
	cmdutil.CONF_SHOW_TOTAL_DURATION: "adds a totals line on time entry " +
		"reports with the sum of the time entries duration",
	cmdutil.CONF_LOG_LEVEL: "how much logs should be shown values: " +
		"none , error , info and debug",
	cmdutil.CONF_ALLOW_ARCHIVED_TAGS: "should allow and suggest archived tags",
	cmdutil.CONF_INTERACTIVE_PAGE_SIZE: "how many entries should be listed " +
		"when prompting options",
	cmdutil.CONF_TIME_ENTRY_DEFAULTS: "should load defaults for time " +
		"entries from current folder",
}

// NewCmdConfig represents the config command
func NewCmdConfig(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manages CLI configuration",
		Args:  cobra.MaximumNArgs(0),
		Example: heredoc.Doc(`
			# cli will guide you to configure the CLI
			$ clockify-cli config init

			# token is the minimum information required for the CLI to work
			$ clockify-cli set token <token>

			# you can see your current parameters using:
			$ clockify-cli get

			# if you wanna see the value of token parameter:
			$ clockify-cli get token
		`),
		Long: heredoc.Doc(`
			Changes or shows configuration settings for clockify-cli

			These are the parameters manageable:
		`) + validParameters.Long(),
	}

	cmd.AddCommand(initialize.NewCmdInit(f))
	cmd.AddCommand(set.NewCmdSet(f, validParameters))
	cmd.AddCommand(get.NewCmdGet(f, validParameters))
	cmd.AddCommand(list.NewCmdList(f))

	return cmd
}
