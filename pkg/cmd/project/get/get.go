package get

import (
	"errors"
	"io"
	"strings"

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

// NewCmdGet looks for a project with the informed ID
func NewCmdGet(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, dto.Project) error,
) *cobra.Command {
	of := util.OutputFlags{}
	p := api.GetProjectParam{}
	cmd := &cobra.Command{
		Use:  "get",
		Args: cmdutil.RequiredNamedArgs("project"),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		Short: "Get a project on a Clockify workspace",
		Example: heredoc.Docf(`
			$ %[1]s 621948458cb9606d934ebb1c
			+--------------------------+-------------------+-----------------------------------------+
			|            ID            |       NAME        |                 CLIENT                  |
			+--------------------------+-------------------+-----------------------------------------+
			| 621948458cb9606d934ebb1c | Clockify Cli      | Special (6202634a28782767054eec26)      |
			+--------------------------+-------------------+-----------------------------------------+

			$ %[1]s cli -q
			621948458cb9606d934ebb1c

			$ %[1]s other --format '{{.Name}} - {{ .Color }} | {{ .ClientID }}'
			Other - #03A9F4 | 6202634a28782767054eec26
		`, "clockify-cli project get"),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if p.ProjectID = strings.TrimSpace(args[0]); p.ProjectID == "" {
				return errors.New("project id should not be empty")
			}

			if err := of.Check(); err != nil {
				return err
			}

			if p.Workspace, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() {
				if p.ProjectID, err = search.GetProjectByName(
					c, f.Config(), p.Workspace, p.ProjectID, ""); err != nil {
					return err
				}
			}

			project, err := c.GetProject(p)
			if err != nil {
				return err
			}
			if project == nil {
				return api.EntityNotFound{
					EntityName: "project",
					ID:         args[0],
				}
			}

			if report != nil {
				return report(cmd.OutOrStdout(), &of, *project)
			}

			return util.ReportOne(*project, cmd.OutOrStdout(), of)
		},
	}

	cmd.Flags().BoolVarP(
		&p.Hydrate, "hydrated", "H", false,
		"projects will have custom fields, tasks and memberships "+
			"filled for json and format outputs")

	util.AddReportFlags(cmd, &of)

	return cmd
}
