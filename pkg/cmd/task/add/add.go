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
package add

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdAdd represents the add command
func NewCmdAdd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a task to the specified project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fl, err := util.TaskReadFlags(cmd, f)
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			project, _ := cmd.Flags().GetString("project")
			if f.Config().IsAllowNameForID() && project != "" {
				if project, err = search.GetProjectByName(
					c, fl.Workspace, project); err != nil {
					return err
				}
			}

			task, err := c.AddTask(api.AddTaskParam{
				Workspace:   fl.Workspace,
				ProjectID:   project,
				Name:        fl.Name,
				Estimate:    fl.Estimate,
				AssigneeIDs: fl.AssigneeIDs,
				Billable:    fl.Billable,
			})
			if err != nil {
				return err
			}

			return util.TaskReport(cmd, task)
		},
	}

	util.TaskAddReportFlags(cmd)

	cmdutil.AddProjectFlags(cmd, f)
	util.TaskAddPropFlags(cmd, f)
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
