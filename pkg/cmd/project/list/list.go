package list

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// projectListCmd represents the projectList command
func NewCmdList(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, []dto.Project) error,
) *cobra.Command {
	of := util.OutputFlags{}
	p := api.GetProjectsParam{
		PaginationParam: api.AllPages(),
	}
	var archived, notArchived bool
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects on a Clockify workspace",
		Example: heredoc.Docf(`
			$ %[1]s
			+--------------------------+-------------------+-----------------------------------------+
			|            ID            |       NAME        |                 CLIENT                  |
			+--------------------------+-------------------+-----------------------------------------+
			| 621948458cb9606d934ebb1c | Clockify Cli      | Special (6202634a28782767054eec26)      |
			| 62a8b52d67f40258719037f2 | New One           |                                         |
			| 62a8b59067f40258719038fc | Other             |                                         |
			| 62a8b607027fe4592ef1520b | Other             | Uber Special (62964b36bb48532a70730dbe) |
			| 62894c3ed2df9d2867dc750b | Something Newer   | Special (6202634a28782767054eec26)      |
			+--------------------------+-------------------+-----------------------------------------+

			$ %[1]s --clients=uber
			+--------------------------+-------------------+-----------------------------------------+
			|            ID            |       NAME        |                 CLIENT                  |
			+--------------------------+-------------------+-----------------------------------------+
			| 62a8b607027fe4592ef1520b | Other             | Uber Special (62964b36bb48532a70730dbe) |
			+--------------------------+-------------------+-----------------------------------------+

			$ %[1]s --clients=uber --clients=special -q
			621948458cb9606d934ebb1c
			62a8b607027fe4592ef1520b

			$ %[1]s --name=other --format '{{.Name}} - {{ .Color }} | {{ .ClientID }}'
			Other - #607D8B | 
			Other - #03A9F4 | 6202634a28782767054eec26

			$ %[1]s --archived
			+--------------------------+-------------------+-----------------------------------------+
			|            ID            |       NAME        |                 CLIENT                  |
			+--------------------------+-------------------+-----------------------------------------+
			| 62894c3ed2df9d2867dc750b | Something Newer   | Special (6202634a28782767054eec26)      |
			+--------------------------+-------------------+-----------------------------------------+
		`, "clockify-cli project list"),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := of.Check(); err != nil {
				return err
			}

			if err := cmdutil.XorFlag(map[string]bool{
				"archived":     archived,
				"not-archived": notArchived,
			}); err != nil {
				return err
			}

			if p.Workspace, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if len(p.Clients) > 0 && f.Config().IsAllowNameForID() {
				if p.Clients, err = search.GetClientsByName(
					c, p.Workspace, p.Clients); err != nil {
					return err
				}
			}

			if archived || notArchived {
				p.Archived = &archived
			}

			projects, err := c.GetProjects(p)
			if err != nil {
				return err
			}

			if report != nil {
				return report(cmd.OutOrStdout(), &of, projects)
			}

			return util.Report(projects, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().StringVarP(&p.Name, "name", "n", "",
		"will be used to filter the project by name")
	cmd.Flags().StringSliceVarP(&p.Clients, "clients", "c", []string{},
		"will be used to filter the project by client id/name")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "clients",
		cmdcomplutil.NewClientAutoComplete(f))

	cmd.Flags().BoolVarP(
		&notArchived, "not-archived", "", false, "list only active projects")
	cmd.Flags().BoolVarP(
		&archived, "archived", "", false, "list only archived projects")
	cmd.Flags().BoolVarP(
		&p.Hydrate, "hydrated", "H", false,
		"projects will have custom fields, tasks and memberships "+
			"filled for json and format outputs")

	util.AddReportFlags(cmd, &of)

	return cmd
}
