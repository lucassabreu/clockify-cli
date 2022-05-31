package user

import (
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
		Aliases: []string{"user"},
		Short:   "List all users on a Workspace",
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
