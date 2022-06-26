package tag

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/tag"
	"github.com/spf13/cobra"
)

// NewCmdTag represents the tags command
func NewCmdTag(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tag",
		Aliases: []string{"tags"},
		Short:   "List tags on Clockify",
		Example: heredoc.Docf(`
			$ %[1]s
			+--------------------------+------------------+
			|            ID            |       NAME       |
			+--------------------------+------------------+
			| 62194867edaba27d0a45b464 | Code Review      |
			| 6219485e8cb9606d934ebb5f | Meeting          |
			| 621948708cb9606d934ebba7 | Pair Programming |
			| 6143b768195e5c503960a775 | Special Tag      |
			+--------------------------+------------------+

			$ %[1]s --name code -q
			62194867edaba27d0a45b464

			$ %[1]s --format "{{.Name}}" -archived
			Archived Tag
		`, "clockify-cli tag"),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			quiet, _ := cmd.Flags().GetBool("quiet")
			if err := cmdutil.XorFlag(map[string]bool{
				"format": format != "",
				"quiet":  quiet,
			}); err != nil {
				return err
			}

			archived, _ := cmd.Flags().GetBool("archived")
			name, _ := cmd.Flags().GetString("name")

			tags, err := getTags(f, name, archived)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if format != "" {
				return output.TagPrintWithTemplate(format)(tags, out)
			}

			if quiet {
				return output.TagPrintQuietly(tags, out)
			}

			return output.TagPrint(tags, os.Stdout)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each Tag")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")
	cmd.Flags().BoolP("archived", "", false, "only display archived tags")

	return cmd
}

func getTags(f cmdutil.Factory, name string, archived bool) ([]dto.Tag, error) {
	c, err := f.Client()
	if err != nil {
		return []dto.Tag{}, err
	}

	w, err := f.GetWorkspaceID()
	if err != nil {
		return []dto.Tag{}, err
	}

	return c.GetTags(api.GetTagsParam{
		Workspace:       w,
		Name:            name,
		Archived:        &archived,
		PaginationParam: api.AllPages(),
	})
}
