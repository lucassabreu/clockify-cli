package version

import (
	"fmt"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdVersion represents the version command
func NewCmdVersion(f cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Shows the CLI version",
		Run: func(cmd *cobra.Command, _ []string) {
			v := f.Version()
			fmt.Fprintln(cmd.OutOrStdout(),
				"Version: "+v.Tag+", Commit: "+v.Commit+", Build At: "+v.Date,
			)
		},
	}
}
