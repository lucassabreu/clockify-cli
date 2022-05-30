package list

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdList represents the list command
func NewCmdList(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tasks of a Clockify project",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			workspace, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			p := api.GetTasksParam{
				Workspace:       workspace,
				PaginationParam: api.AllPages(),
			}

			p.Active, _ = cmd.Flags().GetBool("active")
			p.Name, _ = cmd.Flags().GetString("name")
			p.ProjectID, _ = cmd.Flags().GetString("project")

			if f.Config().IsAllowNameForID() &&
				p.ProjectID != "" {
				if p.ProjectID, err = search.GetProjectByName(
					c, workspace, p.ProjectID); err != nil {
					return err
				}
			}

			tasks, err := c.GetTasks(p)
			if err != nil {
				return err
			}

			return util.TaskReport(cmd, of, tasks...)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().BoolP("active", "a", false, "display only active tasks")

	util.TaskAddReportFlags(cmd, &of)
	cmdutil.AddProjectFlags(cmd, f)

	return cmd
}
