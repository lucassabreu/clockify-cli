package list

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// projectListCmd represents the projectList command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			clients, _ := cmd.Flags().GetStringSlice("clients")

			workspace, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			conf := f.Config()
			if conf.IsAllowNameForID() {
				clients, err = search.GetClientsByName(c, workspace, clients)
				if err != nil {
					return err
				}
			}

			p := api.GetProjectsParam{
				Workspace:       workspace,
				Name:            name,
				Clients:         clients,
				PaginationParam: api.AllPages(),
			}

			if ok, _ := cmd.Flags().GetBool("not-archived"); ok {
				b := false
				p.Archived = &b
			}
			if ok, _ := cmd.Flags().GetBool("archived"); ok {
				b := true
				p.Archived = &b
			}

			projects, err := c.GetProjects(p)
			if err != nil {
				return err
			}

			return util.Report(projects, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the project by name")
	cmd.Flags().StringSliceP("clients", "c", []string{},
		"will be used to filter the project by client id/name")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "clients",
		cmdcomplutil.NewClientAutoComplete(f))

	cmd.Flags().BoolP("not-archived", "", false, "list only active projects")
	cmd.Flags().BoolP("archived", "", false, "list only archived projects")

	util.AddReportFlags(cmd, &of)

	return cmd
}
