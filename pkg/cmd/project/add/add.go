package add

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/project"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdAdd represents the add command
func NewCmdAdd(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new", "create"},
		Short:   "Adds a project to the Clockify workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			note, _ := cmd.Flags().GetString("note")
			color, _ := cmd.Flags().GetString("color")
			randomColor, _ := cmd.Flags().GetBool("random-color")
			client, _ := cmd.Flags().GetString("client")
			public, _ := cmd.Flags().GetBool("public")
			billable, _ := cmd.Flags().GetBool("billable")

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() && client != "" {
				cs, err := search.GetClientsByName(c, w, []string{client})
				if err != nil {
					return err
				}
				client = cs[0]
			}

			if randomColor {
				bytes := make([]byte, 3)
				if _, err := rand.Read(bytes); err != nil {
					return err
				}
				color = "#" + hex.EncodeToString(bytes)
			}

			project, err := c.AddProject(api.AddProjectParam{
				Workspace: w,
				Name:      name,
				ClientId:  client,
				Color:     color,
				Billable:  billable,
				Public:    public,
				Note:      note,
			})
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if of.JSON {
				return output.ProjectJSONPrint(project, out)
			}

			return util.Report([]dto.Project{project}, out, of)
		},
	}

	cmd.Flags().StringP("name", "n", "", "name of the new project")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringP("color", "c", "", "color of the new project")
	cmd.Flags().Bool("random-color", false, "use a random color for the project")
	cmd.Flags().StringP("note", "N", "", "note for the new project")
	cmd.Flags().String("client", "", "the id/name of the client the new project will go under")
	cmd.Flags().BoolP("public", "p", false, "make the new project public")
	cmd.Flags().BoolP("billable", "b", false, "make the new project as billable")

	util.AddReportFlags(cmd, &of)

	return cmd
}
