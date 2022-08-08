package edit

import (
	"errors"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdEdit represents the close command
func NewCmdEdit(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, dto.Task) error,
) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "edit <task>",
		Aliases: []string{"update"},
		Args:    cmdutil.RequiredNamedArgs("task"),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewTaskAutoComplete(f, false)),
		Short: "Edit a task from a project on Clockify",
		Long: heredoc.Doc(`
			Edits a task on a Clockify's project, allowing to change the name, estimated time, assignees, status and billable settings.

			If you set a estimate for the task, but the project is set as manual estimation, then it will have no effect on Clockify.
		`),
		Example: heredoc.Docf(`
			$ %[1]s -p special 62aa5d7049445270d7b979d6 --name="Very Important"
			+--------------------------+----------------+--------+
			|            ID            |      NAME      | STATUS |
			+--------------------------+----------------+--------+
			| 62aa5d7049445270d7b979d6 | Very Important | ACTIVE |
			+--------------------------+----------------+--------+

			$ %[1]s -p special 'important' --assign john@example.com | \
			  jq '.[] |.assigneeIds' --compact-output
			["dddddddddddddddddddddddd"]

			$ %[1]s -p special important --billable --quiet
			62aa5d7049445270d7b979d6

			$ %[1]s -p special important --not-billable --csv
			id,name,status
			62aa5d7049445270d7b979d6,Very Important,ACTIVE

			$ %[1]s -p special very --estimate 1 --done
			+--------------------------+----------------+--------+
			|            ID            |      NAME      | STATUS |
			+--------------------------+----------------+--------+
			| 62aa5d7049445270d7b979d6 | Very Important | DONE   |
			+--------------------------+----------------+--------+

			$ %[1]s -p special 'very i' --active --format --no-assignee \
			  --format '{{.Name}} | {{.Status}} | {{ .AssigneeIDs }}'
			Very Important | ACTIVE | []
		`, "clockify-cli task edit"),
		RunE: func(cmd *cobra.Command, args []string) error {
			task := strings.TrimSpace(args[0])
			if task == "" {
				return errors.New("task id should not be empty")
			}

			if err := of.Check(); err != nil {
				return err
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
				if task, err = search.GetTaskByName(
					c,
					api.GetTasksParam{
						Workspace: fl.Workspace, ProjectID: fl.ProjectID},
					task); err != nil {
					return err
				}
			}

			p := api.UpdateTaskParam{
				Workspace:   fl.Workspace,
				ProjectID:   fl.ProjectID,
				TaskID:      task,
				Name:        fl.Name,
				Estimate:    fl.Estimate,
				AssigneeIDs: fl.AssigneeIDs,
				Billable:    fl.Billable,
			}

			if !cmd.Flags().Changed("name") {
				t, err := c.GetTask(api.GetTaskParam{
					Workspace: fl.Workspace,
					ProjectID: fl.ProjectID,
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

			if report == nil {
				return util.TaskReport(cmd, of, t)
			}

			return report(cmd.OutOrStdout(), &of, t)
		},
	}

	util.TaskAddPropFlags(cmd, f)

	cmd.Flags().Bool("done", false, "sets the task as done")
	cmd.Flags().Bool("active", false, "sets the task as active")

	util.TaskAddReportFlags(cmd, &of)

	return cmd
}
