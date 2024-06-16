package show

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdShow prints the default options for the current folder
func NewCmdShow(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}
