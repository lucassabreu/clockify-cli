package list

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdList represents the list command
func NewCmdList(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, []dto.Task) error,
) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks in a Clockify project",
		Example: heredoc.Docf(`
			$ %[1]s --project special
			+--------------------------+----------+--------+
			|            ID            |   NAME   | STATUS |
			+--------------------------+----------+--------+
			| 62aa4eed49445270d7b9666c | Inactive | DONE   |
			| 62aa4ee64ebb4f143c8d5225 | Second   | ACTIVE |
			| 62aa4ea2c22de9759e6e3a0e | First    | ACTIVE |
			+--------------------------+----------+--------+

			$ %[1]s --project special --active --quiet
			62aa4ee64ebb4f143c8d5225
			62aa4ea2c22de9759e6e3a0e

			$ %[1]s --project special --name inact --csv
			id,name,status
			62aa4eed49445270d7b9666c,Inactive,DONE
		`, "clockify-cli task list"),
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

			if report == nil {
				return util.TaskReport(cmd, of, tasks...)
			}

			return report(cmd.OutOrStdout(), &of, tasks)
		},
	}

	cmd.Flags().StringP("name", "n", "",
		"will be used to filter the tag by name")
	cmd.Flags().BoolP("active", "a", false, "display only active tasks")

	util.TaskAddReportFlags(cmd, &of)
	cmdutil.AddProjectFlags(cmd, f)

	return cmd
}
