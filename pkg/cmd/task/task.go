package task

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/add"
	del "github.com/lucassabreu/clockify-cli/pkg/cmd/task/delete"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/done"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/edit"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdTask represents the client command
func NewCmdTask(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "task",
		Aliases: []string{"tasks"},
		Short:   "Work with Clockify tasks",
	}

	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(add.NewCmdAdd(f))
	cmd.AddCommand(edit.NewCmdEdit(f))
	cmd.AddCommand(del.NewCmdDelete(f))
	cmd.AddCommand(done.NewCmdDone(f))

	return cmd
}
