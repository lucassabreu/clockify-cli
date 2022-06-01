package edit

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

// NewCmdEdit represents the close command
func NewCmdEdit(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "edit <project> <task>",
		Aliases: []string{"update"},
		Args:    cobra.ExactArgs(2),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		Short: "Edit a task from a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New(
					"two arguments are required (project and task)")
			}

			project := strings.TrimSpace(args[0])
			task := strings.TrimSpace(args[1])
			if project == "" || task == "" {
				return errors.New("project and task id should not be empty")
			}

			fl, err := util.TaskReadFlags(cmd, f)
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() {
				if project, err = search.GetProjectByName(
					c, fl.Workspace, project); err != nil {
					return err
				}

				if task, err = search.GetTaskByName(
					c, fl.Workspace, project, task); err != nil {
					return err
				}

			}

			p := api.UpdateTaskParam{
				Workspace:   fl.Workspace,
				ProjectID:   project,
				TaskID:      task,
				Name:        fl.Name,
				Estimate:    fl.Estimate,
				AssigneeIDs: fl.AssigneeIDs,
				Billable:    fl.Billable,
			}

			if !cmd.Flags().Changed("name") {
				t, err := c.GetTask(api.GetTaskParam{
					Workspace: fl.Workspace,
					ProjectID: project,
					TaskID:    task,
				})
				if err != nil {
					return err
				}

				p.Name = t.Name
			}

			switch {
			case cmd.Flags().Changed("active"):
				p.Status = api.TaskStatusActive
			case cmd.Flags().Changed("done"):
				p.Status = api.TaskStatusDone
			default:
				p.Status = api.TaskStatusDefault
			}

			t, err := c.UpdateTask(p)
			if err != nil {
				return err
			}

			return util.TaskReport(cmd, of, t)
		},
	}

	util.TaskAddPropFlags(cmd, f)
	cmdutil.AddProjectFlags(cmd, f)

	cmd.Flags().Bool("done", false, "sets the task as done")
	cmd.Flags().Bool("active", false, "sets the task as active")

	util.TaskAddReportFlags(cmd, &of)

	return cmd
}
