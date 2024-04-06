package defaults

import (
	"io"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/show"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	outd "github.com/lucassabreu/clockify-cli/pkg/output/defaults"
	"github.com/spf13/cobra"
)

// NewCmdDefaults creates commands to manage default parameters when creating
// time entries
func NewCmdDefaults(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "defaults",
		Aliases: []string{"def"},
		Short: "Manages the default parameters for time entries " +
			"in the current folder",
		Args: cobra.ExactArgs(0),
	}

	cmd.AddCommand(
		set.NewCmdSet(f, func(of outd.OutputFlags, w io.Writer, dte defaults.DefaultTimeEntry) error {
			return nil
		}),
		show.NewCmdShow(f),
	)

	return cmd
}
