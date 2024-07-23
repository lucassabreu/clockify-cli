package show

import (
	"io"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	outd "github.com/lucassabreu/clockify-cli/pkg/output/defaults"
	"github.com/spf13/cobra"
)

// NewCmdShow prints the default options for the current folder
func NewCmdShow(
	f cmdutil.Factory,
	report func(outd.OutputFlags, io.Writer, defaults.DefaultTimeEntry) error,
) *cobra.Command {
	of := outd.OutputFlags{}
	cmd := &cobra.Command{
		Use:  "show",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := f.TimeEntryDefaults().Read()
			if err != nil {
				return err
			}

			return report(of, cmd.OutOrStdout(), d)
		},
	}

	cmd.Flags().StringVarP(&of.Format,
		"format", "f", outd.FORMAT_YAML, "output format")
	_ = cmdcompl.AddFixedSuggestionsToFlag(cmd, "format",
		cmdcompl.ValidArgsSlide{outd.FORMAT_YAML, outd.FORMAT_JSON})

	return cmd
}
