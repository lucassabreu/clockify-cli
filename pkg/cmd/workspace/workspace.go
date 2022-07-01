package workspace

import (
	"os"

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
		Short:   "List your available workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			quiet, _ := cmd.Flags().GetBool("quiet")

			if err := cmdutil.XorFlag(map[string]bool{
				"format": format != "",
				"quiet":  quiet,
			}); err != nil {
				return err
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			list, err := c.GetWorkspaces(api.GetWorkspaces{
				Name: name,
			})
			if err != nil {
				return err
			}

			w, _ := f.GetWorkspaceID()
			if quiet {
				return output.WorkspacePrintQuietly(list, os.Stdout)
			}

			if format != "" {
				return output.WorkspacePrintWithTemplate(format)(
					list, os.Stdout)
			}

			return output.WorkspacePrint(w)(list, os.Stdout)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the workspaces by name")
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each workspace")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}
