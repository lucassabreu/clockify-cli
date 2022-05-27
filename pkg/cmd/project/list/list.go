package list

import (
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
		Short:   "List projects on Clockify and project links",
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
