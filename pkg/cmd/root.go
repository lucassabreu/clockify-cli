package cmd

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/completion"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/tag"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task"
	timeentry "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry"
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

	cmd.PersistentFlags().BoolP("interactive", "i", false,
		"will prompt you to confirm/complement commands input before "+
			"executing the action ")

	cmd.PersistentFlags().BoolP("allow-name-for-id", "", false,
		"allow use of project/client/tag's name when id is asked")

	cmd.PersistentFlags().String(
		"log-level", cmdutil.LOG_LEVEL_NONE, "set log level")
	_ = cmdcompl.AddFixedSuggestionsToFlag(cmd, "log-level",
		cmdcompl.ValidArgsSlide{
			cmdutil.LOG_LEVEL_NONE,
			cmdutil.LOG_LEVEL_DEBUG,
			cmdutil.LOG_LEVEL_INFO,
		})

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

	return cmd
}
