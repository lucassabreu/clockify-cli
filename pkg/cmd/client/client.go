package client

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/add"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdClient represents the client command
func NewCmdClient(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"clients"},
		Short:   "Work with Clockify clients",
	}

	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(add.NewCmdAdd(f))

	return cmd
}
