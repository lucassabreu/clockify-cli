package add

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/output/client"
	"github.com/spf13/cobra"
)

// clientAddCmd represents the add command
func NewCmdAdd(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "create"},
		Short:   "Adds a client to the Clockify workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			cl, err := c.AddClient(api.AddClientParam{
				Workspace: w,
				Name:      name,
			})
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if of.JSON {
				client.ClientJSONPrint(cl, out)
			}

			return util.Report([]dto.Client{cl}, out, of)
		},
	}

	cmd.Flags().StringP("name", "n", "", "the name of the new client")
	_ = cmd.MarkFlagRequired("name")

	util.AddReportFlags(cmd, &of)

	return cmd
}
