package list

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList represents the list command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List clients on Clockify",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			p := api.GetClientsParam{
				PaginationParam: api.AllPages(),
			}

			var err error
			if p.Workspace, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			p.Name, _ = cmd.Flags().GetString("name")
			if ok, _ := cmd.Flags().GetBool("not-archived"); ok {
				b := false
				p.Archived = &b
			}
			if ok, _ := cmd.Flags().GetBool("archived"); ok {
				b := true
				p.Archived = &b
			}

			clients, err := c.GetClients(p)
			if err != nil {
				return err
			}

			return util.Report(clients, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().BoolP("not-archived", "", false, "list only active projects")
	cmd.Flags().BoolP("archived", "", false, "list only archived projects")

	util.AddReportFlags(cmd, &of)

	return cmd
}
