package workspace

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/workspace"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// NewCmdWorkspace represents the workspaces command
func NewCmdWorkspace(f cmdutil.Factory) *cobra.Command {
	fl := struct {
		name   string
		format string
		quiet  bool
	}{}
	cmd := &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"workspaces"},
		Short:   "List your available workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.XorFlag(map[string]bool{
				"format": fl.format != "",
				"quiet":  fl.quiet,
			}); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			list, err := c.GetWorkspaces(api.GetWorkspaces{
				Name: fl.name,
			})
			if err != nil {
				return err
			}

			if fl.quiet {
				return output.WorkspacePrintQuietly(list, cmd.OutOrStdout())
			}

			if fl.format != "" {
				return output.WorkspacePrintWithTemplate(fl.format)(
					list, cmd.OutOrStdout())
			}

			w, _ := f.GetWorkspaceID()
			return output.WorkspacePrint(w)(list, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&fl.name, "name", "n", "",
		"will be used to filter the workspaces by name")
	cmd.Flags().StringVarP(&fl.format, "format", "f", "",
		"golang text/template format to be applied on each workspace")
	cmd.Flags().BoolVarP(&fl.quiet, "quiet", "q", false, "only display ids")

	return cmd
}
