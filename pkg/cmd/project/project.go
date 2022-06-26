package project

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/add"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdProject represents the project command
func NewCmdProject(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects"},
		Short:   "Work with Clockify projects",
	}

	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(add.NewCmdAdd(f))

	return cmd
}
