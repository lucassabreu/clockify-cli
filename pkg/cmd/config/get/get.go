package get

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdGet(
	f cmdutil.Factory,
	validParameters cmdcompl.ValidArgsMap,
) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get <param>",
		Short: "Retrieves one parameter set by the user",
		Example: heredoc.Docf(`
			$ %[1]s token
			Yamdas569

			$ %[1]s workweek-days --format=json
			["monday","tuesday","wednesday","thursday","friday"]
		`, "clockify-cli config get"),
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("param"),
			cobra.ExactArgs(1),
		),
		ValidArgs: validParameters.IntoValidArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return util.Report(
				cmd.OutOrStdout(), format,
				f.Config().Get(args[0]))
		},
	}

	_ = util.AddReportFlags(cmd, &format)

	return cmd
}
