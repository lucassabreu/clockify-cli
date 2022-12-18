package task

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/add"
	del "github.com/lucassabreu/clockify-cli/pkg/cmd/task/delete"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/done"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/edit"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/list"
	quickadd "github.com/lucassabreu/clockify-cli/pkg/cmd/task/quick-add"
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

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(add.NewCmdAdd(f, nil))
	cmd.AddCommand(quickadd.NewCmdQuickAdd(f, nil))
	cmd.AddCommand(edit.NewCmdEdit(f, nil))
	cmd.AddCommand(del.NewCmdDelete(f, nil))
	cmd.AddCommand(done.NewCmdDone(f, nil))

	return cmd
}
