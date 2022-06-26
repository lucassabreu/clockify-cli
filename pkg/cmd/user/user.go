package user

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/me"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUser represents the users command
func NewCmdUser(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		Short:   "List users of a workspace",
		Example: heredoc.Docf(`
			$ %[1]s
			+--------------------------+-------------+--------------+--------+
			|            ID            |    NAME     |     EMAIL    | STATUS |
			+--------------------------+-------------+--------------+--------+
			| eeeeeeeeeeeeeeeeeeeeeeee | John Due    | john@due.net | ACTIVE |
			| ffffffffffffffffffffffff | John JD Due | due@john.net | ACTIVE |
			+--------------------------+-------------+--------------+--------+

			$ %[1]s --quiet
			eeeeeeeeeeeeeeeeeeeeeeee
			ffffffffffffffffffffffff

			$ %[1]s --email due@john.net
			+--------------------------+-------------+--------------+--------+
			|            ID            |    NAME     |     EMAIL    | STATUS |
			+--------------------------+-------------+--------------+--------+
			| ffffffffffffffffffffffff | John JD Due | due@john.net | ACTIVE |
			+--------------------------+-------------+--------------+--------+

			$ %[1]s me --format "{{ .Name }} ({{ .Email }})" --email due@john.net
			John JD Due (due@john.net)
		`, "clockify-cli user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			email, _ := cmd.Flags().GetString("email")
			if err := of.Check(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
				Workspace: w,
				Email:     email,
			})
			if err != nil {
				return err
			}

			return util.Report(users, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().StringP("email", "e", "",
		"will be used to filter the workspaces by email")

	util.AddReportFlags(cmd, &of)

	_ = cmd.MarkFlagRequired("workspace")

	cmd.AddCommand(me.NewCmdMe(f))

	return cmd
}
