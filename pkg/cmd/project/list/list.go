package list

import (
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/project"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// projectListCmd represents the projectList command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects on Clockify and project links",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			asJSON, _ := cmd.Flags().GetBool("json")
			asCSV, _ := cmd.Flags().GetBool("csv")
			quiet, _ := cmd.Flags().GetBool("quiet")
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

			return report(projects, cmd.OutOrStdout(),
				format, asJSON, asCSV, quiet)
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
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each Project")

	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}

func report(projects []dto.Project, out io.Writer,
	format string, asJSON, asCSV, quiet bool) error {

	if asJSON {
		return output.ProjectsJSONPrint(projects, out)
	}

	if asCSV {
		return output.ProjectsCSVPrint(projects, out)
	}

	if format != "" {
		return output.ProjectPrintWithTemplate(format)(projects, out)
	}

	if quiet {
		return output.ProjectPrintQuietly(projects, out)
	}

	return output.ProjectPrint(projects, os.Stdout)

}
