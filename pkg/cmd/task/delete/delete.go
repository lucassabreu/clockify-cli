package del

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// NewCmdDelete represents the close command
func NewCmdDelete(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "delete <task>",
		Aliases: []string{"remove", "rm", "del"},
		Args:    cmdutil.RequiredNamedArgs("task"),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewTaskAutoComplete(f, false)),
		Short: "Deletes a task from a project on Clockify",
		Long: heredoc.Doc(`
			Deletes a task from a project on Clockify
			This action can't be reverted, and all time entries using this task will revert to not having one
		`),
		Example: heredoc.Doc(`
			$ clockify-cli task delete -p "special" very
			+--------------------------+----------------+--------+
			|            ID            |      NAME      | STATUS |
			+--------------------------+----------------+--------+
			| 62aa5d7049445270d7b979d6 | Very Important | ACTIVE |
			+--------------------------+----------------+--------+

			$ clockify-cli task delete -p "special" 62aa4eed49445270d7b9666c
			+--------------------------+----------+--------+
			|            ID            |   NAME   | STATUS |
			+--------------------------+----------+--------+
			| 62aa4eed49445270d7b9666c | Inactive | DONE   |
			+--------------------------+----------+--------+

			$ clockify-cli task delete -p "special" 62aa4eed49445270d7b9666c
			No task with id or name containing '62aa4eed49445270d7b9666c' was found
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, _ := cmd.Flags().GetString("project")
			task := strings.TrimSpace(args[0])
			if project == "" || task == "" {
				return errors.New("project and task id should not be empty")
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if f.Config().IsAllowNameForID() {
				if project, err = search.GetProjectByName(
					c, w, project); err != nil {
					return err
				}

				if task, err = search.GetTaskByName(
					c,
					api.GetTasksParam{Workspace: w, ProjectID: project},
					task,
				); err != nil {
					return err
				}
			}

			t, err := c.DeleteTask(api.DeleteTaskParam{
				Workspace: w,
				ProjectID: project,
				TaskID:    task,
			})
			if err != nil {
				return err
			}

			return util.TaskReport(cmd, of, t)
		},
	}

	cmdutil.AddProjectFlags(cmd, f)
	util.TaskAddReportFlags(cmd, &of)

	return cmd
}
