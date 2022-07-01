package me

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/output/user"
	"github.com/spf13/cobra"
)

// NewCmdMe represents the me command
func NewCmdMe(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Show details about the user who owns the token",
		Long: heredoc.Doc(`
			Shows details about the user who owns the token used by the CLI.

			This user may be different from the one set at "user.id", but if the parameter is not set the CLI will defaults to this one.
		`),
		Example: heredoc.Docf(`
			$ %[1]s
			+--------------------------+-------------+--------------+--------+
			|            ID            |    NAME     |     EMAIL    | STATUS |
			+--------------------------+-------------+--------------+--------+
			| ffffffffffffffffffffffff | John JD Due | due@john.net | ACTIVE |
			+--------------------------+-------------+--------------+--------+

			$ %[1]s --quiet
			ffffffffffffffffffffffff

			$ %[1]s --format "{{ .Name }} ({{ .Email }})"
			John JD Due (due@john.net)
		`, "clockify-cli user me"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			u, err := c.GetMe()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if of.JSON {
				return user.UserJSONPrint(u, out)
			}

			return util.Report([]dto.User{u}, out, of)
		},
	}

	util.AddReportFlags(cmd, &of)

	return cmd
}
