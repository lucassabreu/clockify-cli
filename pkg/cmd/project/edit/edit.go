package edit

import (
	"errors"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// NewCmdEdit updates a project
func NewCmdEdit(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, []dto.Project) error,
) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "edit <project>...",
		Aliases: []string{"update"},
		Short:   "Edit a project",
		Example: heredoc.Docf(`
			# set a client form the project
			$ clockify-cli project edit cli --client Myself
			+--------------------------+--------------+--------------------------------+
			|            ID            |     NAME     |             CLIENT             |
			+--------------------------+--------------+--------------------------------+
			| 621948458cb9606d934ebb1c | Clockify Cli | Myself                         |
			|                          |              | (6202634a28782767054eec26)     |
			+--------------------------+--------------+--------------------------------+

			# remove client from a project
			$ clockify-cli project edit cli --no-client
			+--------------------------+--------------+--------+
			|            ID            |     NAME     | CLIENT |
			+--------------------------+--------------+--------+
			| 621948458cb9606d934ebb1c | Clockify Cli |        |
			+--------------------------+--------------+--------+

			# change name, color and make public
			$ clockify-cli project 62f19c254a912b05acc6d6cf \
				--name First --public --color #0f0 \
				--format "{{.Name}} | {{.Public}} | {{.Color}}"
			First | true | #00ff00

			# change to not billable, archived and leave a note
			$ clockify-cli project second --not-billable --archived \
				--note "$(cat notes.txt)" \
				--format 'n: {{.Name}}\nb: {{.Billable}}\na: {{.Archived}}\nn:\n{{ .Note }}'
			n: Noted
			b: false
			a: false
			n: one line
			two lines
			three lines

			# archive multiple projects
			$ clockify-cli project first second \
				--archived \
				--format "{{.Name}} | {{.Archived}}"
			First | true
			Second | true
		`),
		Args: cmdutil.RequiredNamedArgs("project"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			if err := cmdutil.XorFlagSet(
				cmd.Flags(), "billable", "not-billable"); err != nil {
				return err
			}

			if err := cmdutil.XorFlagSet(
				cmd.Flags(), "private", "public"); err != nil {
				return err
			}

			if err := cmdutil.XorFlagSet(
				cmd.Flags(), "no-client", "client"); err != nil {
				return err
			}

			if err := cmdutil.XorFlagSet(
				cmd.Flags(), "archived", "active"); err != nil {
				return err
			}

			if len(args) > 1 && cmd.Flags().Changed("name") {
				return errors.New(
					"`--name` can't be changed for multiple projects")
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			ids := strhlp.Unique(strhlp.Map(strings.TrimSpace, args))
			var client *string

			if cmd.Flags().Changed("client") {
				id, _ := cmd.Flags().GetString("client")
				client = &id
			} else if cmd.Flags().Changed("no-client") {
				id := ""
				client = &id
			}

			if f.Config().IsAllowNameForID() {
				if ids, err = search.GetProjectsByName(
					c, f.Config(), w, "", ids); err != nil {
					return err
				}

				if client != nil && *client != "" {
					if *client, err = search.GetClientByName(
						c, w, *client); err != nil {
						return err
					}
				}
			}

			p := api.UpdateProjectParam{
				Workspace: w,
				ClientId:  client,
			}

			p.Name, _ = cmd.Flags().GetString("name")
			p.Color, _ = cmd.Flags().GetString("color")

			if cmd.Flags().Changed("note") {
				n, _ := cmd.Flags().GetString("note")
				p.Note = &n
			}

			if cmd.Flags().Changed("billable") ||
				cmd.Flags().Changed("not-billable") {
				b, _ := cmd.Flags().GetBool("billable")
				p.Billable = &b
			}

			if cmd.Flags().Changed("archived") ||
				cmd.Flags().Changed("active") {
				b, _ := cmd.Flags().GetBool("archived")
				p.Archived = &b
			}

			if cmd.Flags().Changed("public") ||
				cmd.Flags().Changed("private") {
				b, _ := cmd.Flags().GetBool("public")
				p.Public = &b
			}

			var g errgroup.Group
			projects := make([]dto.Project, len(ids))
			for i := 0; i < len(ids); i++ {
				j := i
				g.Go(func() error {
					cp := p
					cp.ProjectID = ids[j]
					projects[j], err = c.UpdateProject(cp)
					return err
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			if report == nil {
				return util.Report(projects, cmd.OutOrStdout(), of)
			}

			return report(cmd.OutOrStdout(), &of, projects)
		},
	}
	cmd.Flags().StringP("name", "n", "", "name of the project")

	cmd.Flags().StringP("color", "c", "",
		"color of the projects")
	cmd.Flags().StringP("note", "N", "",
		"note for the projects")

	cmd.Flags().String("client", "",
		"the id/name of the client the projects will go under")
	cmd.Flags().Bool("no-client", false,
		"set projects as not having clients")

	cmd.Flags().BoolP("public", "p", false,
		"set projects as public")
	cmd.Flags().BoolP("private", "P", false,
		"set the projects as private")

	cmd.Flags().BoolP("billable", "b", false,
		"set the projects as billable")
	cmd.Flags().BoolP("not-billable", "B", false,
		"set the projects as not billable")

	cmd.Flags().BoolP("archived", "A", false,
		"set projects as archived")
	cmd.Flags().BoolP("active", "a", false,
		"set the projects as active")

	util.AddReportFlags(cmd, &of)

	return cmd
}
