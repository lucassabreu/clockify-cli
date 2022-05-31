package me

import (
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
		Short: "Show the user's token info",
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
