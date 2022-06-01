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
	of := util.OutputFlags{}
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

			return util.TaskReport(cmd, of, task)
		},
	}

	util.TaskAddReportFlags(cmd, &of)

	cmdutil.AddProjectFlags(cmd, f)
	util.TaskAddPropFlags(cmd, f)
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
