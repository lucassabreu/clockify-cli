package list

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the config list command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all parameters set by the user",
		Example: heredoc.Doc(`
			$ clockify-cli config list
			allow-incomplete: false
			allow-name-for-id: true
			allow-project-name: true
			debug: false
			description-autocomplete: true
			description-autocomplete-days: 15
			interactive: true
			no-closing: false
			show-task: false
			show-total-duration: true
			token: Yamdas569
			user:
			  id: ffffffffffffffffffffffff
			workspace: eeeeeeeeeeeeeeeeeeeeeeee
			workweek-days:
			- monday
			- tuesday
			- wednesday
			- thursday
			- friday
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return util.Report(cmd.OutOrStdout(), format, f.Config().All())
		},
	}

	_ = util.AddReportFlags(cmd, &format)

	return cmd
}
