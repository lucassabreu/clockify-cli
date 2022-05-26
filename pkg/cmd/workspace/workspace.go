package workspace

import (
	"errors"
	"os"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/workspace"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// NewCmdWorkspace represents the workspaces command
func NewCmdWorkspace(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"workspaces"},
		Short:   "List user's workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("format") && cmd.Flags().Changed("quiet") {
				return errors.New(
					"format and quiet flags can't be used together")
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			list, err := run(c, name)
			if err != nil {
				return err
			}

			w, _ := f.GetWorkspaceID()
			format, _ := cmd.Flags().GetString("format")
			quiet, _ := cmd.Flags().GetBool("quiet")
			return report(list, w, format, quiet)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the workspaces by name")
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each workspace")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}

func run(c *api.Client, name string) ([]dto.Workspace, error) {
	return c.GetWorkspaces(api.GetWorkspaces{
		Name: name,
	})
}

func report(
	list []dto.Workspace, dWorkspace, format string, quiet bool) error {
	if quiet {
		return output.WorkspacePrintQuietly(list, os.Stdout)
	}

	if format != "" {
		return output.WorkspacePrintWithTemplate(format)(list, os.Stdout)
	}

	return output.WorkspacePrint(dWorkspace)(list, os.Stdout)
}
