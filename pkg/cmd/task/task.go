/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
		Short:   "List/add tasks of/to a project",
	}

	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(add.NewCmdAdd(f))
	cmd.AddCommand(edit.NewCmdEdit(f))
	cmd.AddCommand(del.NewCmdDelete(f))
	cmd.AddCommand(done.NewCmdDone(f))

	return cmd
}
