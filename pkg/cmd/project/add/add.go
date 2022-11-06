package add

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdAdd represents the add command
func NewCmdAdd(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, dto.Project) error,
) *cobra.Command {
	of := util.OutputFlags{}
	p := api.AddProjectParam{}
	randomColor := false
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "create"},
		Short:   "Adds a project to the Clockify workspace",
		Example: heredoc.Docf(`
			$ %[1]s --name "New One"
			+--------------------------+---------+--------+
			|            ID            |  NAME   | CLIENT |
			+--------------------------+---------+--------+
			| 62a8b52d67f40258719037f2 | New One |        |
			+--------------------------+---------+--------+

			$ %[1]s --name=Other -q
			62a8b59067f40258719038fc

			$ %[1]s --name "Other" --client="Uber" --csv --color=#fff
			id,name,client.id,client.name
			62a8b607027fe4592ef1520b,Other,62964b36bb48532a70730dbe,Uber Special

			$ %[1]s --name Other --random-color
			add project: Other project for client Uber Special already exists. (code: 501)

			$ %[1]s --name "Something" --client="Uber" --color=#fff
			the following flags can't be used together: color and random-color

			$ %[1]s --name "Something" --client="Uber"
			the following flags can't be used together: color and random-color

		`, "clockify-cli project add"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			if err := cmdutil.XorFlag(map[string]bool{
				"color":        p.Color != "",
				"random-color": randomColor,
			}); err != nil {
				return err
			}

			var err error
			p.Workspace, err = f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if p.ClientId != "" && f.Config().IsAllowNameForID() {
				cs, err := search.GetClientsByName(
					c, p.Workspace, []string{p.ClientId})
				if err != nil {
					return err
				}
				p.ClientId = cs[0]
			}

			if randomColor {
				bytes := make([]byte, 3)
				if _, err := rand.Read(bytes); err != nil {
					return err
				}
				p.Color = "#" + hex.EncodeToString(bytes)
			}

			project, err := c.AddProject(p)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if report != nil {
				return report(out, &of, project)
			}

			return util.ReportOne(project, out, of)
		},
	}

	cmd.Flags().StringVarP(&p.Name, "name", "n", "", "name of the new project")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringVarP(&p.Color, "color", "c", "",
		"color of the new project")
	cmd.Flags().BoolVar(&randomColor, "random-color", false,
		"use a random color for the project")
	cmd.Flags().StringVarP(&p.Note, "note", "N", "",
		"note for the new project")
	cmd.Flags().StringVar(&p.ClientId, "client", "",
		"the id/name of the client the new project will go under")
	cmd.Flags().BoolVarP(&p.Public, "public", "p", false,
		"make the new project public")
	cmd.Flags().BoolVarP(&p.Billable, "billable", "b", false,
		"make the new project as billable")

	util.AddReportFlags(cmd, &of)

	return cmd
}
