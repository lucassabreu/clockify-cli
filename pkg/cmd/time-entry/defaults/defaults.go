package defaults

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/show"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
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
		set.NewCmdSet(f, nil),
		show.NewCmdShow(f),
	)

	return cmd
}
