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
package done

import (
	"errors"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdDone represents the close command
func NewCmdDone(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "done <project> <task>",
		Aliases: []string{"mark-as-done", "end"},
		Args:    cobra.ExactArgs(2),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		Short: "Change a task from a project to done",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("two arguments are required (project and task)")
			}

			project := strings.TrimSpace(args[0])
			task := strings.TrimSpace(args[1])
			if project == "" || task == "" {
				return errors.New("project and task id should not be empty")
			}

			workspace, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() {
				if project, err = search.GetProjectByName(
					c, workspace, project); err != nil {
					return err
				}

				if task, err = search.GetTaskByName(
					c, workspace, project, task); err != nil {
					return err
				}
			}

			t, err := c.GetTask(api.GetTaskParam{
				Workspace: workspace,
				ProjectID: project,
				TaskID:    task,
			})
			if err != nil {
				return err
			}

			if t, err = c.UpdateTask(api.UpdateTaskParam{
				Workspace: workspace,
				ProjectID: project,
				TaskID:    task,
				Name:      t.Name,
				Status:    api.TaskStatusDone,
			}); err != nil {
				return err
			}

			return util.TaskReport(cmd, t)
		},
	}

	util.TaskAddReportFlags(cmd)
	return cmd
}
