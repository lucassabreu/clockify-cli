package list

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList represents the list command
func NewCmdList(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, []dto.Client) error,
) *cobra.Command {
	of := util.OutputFlags{}
	var archived, notArchived bool
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List clients from a Clockify workspace",
		Example: heredoc.Docf(`
			$ %[1]s
			+--------------------------+----------+----------+
			|            ID            |   NAME   | ARCHIVED |
			+--------------------------+----------+----------+
			| 6202634a28782767054eec26 | Client 1 | NO       |
			| 62964b36bb48532a70730dbe | Client 2 | YES      |
			+--------------------------+----------+----------+

			$ %[1]s --archived --csv
			62964b36bb48532a70730dbe,Client 2,true

			$ %[1]s --not-archived --format "<{{ .Name }}>"
			<Client 1>

			$ %[1]s --name "1" --quiet
			6202634a28782767054eec26
		`, "clockify-cli client list"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			if err := cmdutil.XorFlag(map[string]bool{
				"archived":     archived,
				"not-archived": notArchived,
			}); err != nil {
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
			if archived || notArchived {
				p.Archived = &archived
			}

			clients, err := c.GetClients(p)
			if err != nil {
				return err
			}

			if report != nil {
				return report(cmd.OutOrStdout(), &of, clients)
			}

			return util.Report(clients, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().BoolVarP(
		&notArchived, "not-archived", "", false, "list only active projects")
	cmd.Flags().BoolVarP(
		&archived, "archived", "", false, "list only archived projects")

	util.AddReportFlags(cmd, &of)

	return cmd
}
