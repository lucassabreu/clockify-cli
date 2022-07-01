package done

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// NewCmdDone represents the close command
func NewCmdDone(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "done <task>...",
		Aliases: []string{"mark-as-done", "end"},
		Args:    cmdutil.RequiredNamedArgs("task"),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewTaskAutoComplete(f, true)),
		Short: "Edits a task  to done",
		Long:  "Edits a task to done, similar to doing `task edit <task> --done`",
		Example: heredoc.Docf(`
			$ %[1]s ls
			+--------------------------+--------+--------+
			|            ID            |  NAME  | STATUS |
			+--------------------------+--------+--------+
			| 62adfcc8c22de9759e739d66 | Five   | ACTIVE |
			| 62adfcc4c22de9759e739d64 | Four   | ACTIVE |
			| 62adfcb649445270d7becfca | Three  | ACTIVE |
			| 62adfcb149445270d7becfc8 | Second | ACTIVE |
			| 62adfcaa4ebb4f143c92bf8b | First  | ACTIVE |
			+--------------------------+--------+--------+

			$ %[1]s done first second 62adfcb649445270d7becfca
			+--------------------------+--------+--------+
			|            ID            |  NAME  | STATUS |
			+--------------------------+--------+--------+
			| 62adfcaa4ebb4f143c92bf8b | First  | DONE   |
			| 62adfcb149445270d7becfc8 | Second | DONE   |
			| 62adfcb649445270d7becfca | Three  | DONE   |
			+--------------------------+--------+--------+

			$ %[1]s done four
			id,name,status
			62adfcc4c22de9759e739d64,Four,DONE

			$ %[1]s done five
			No active task with id or name containing 'five' was found
		`, "clockify-cli task -p cli"),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, _ := cmd.Flags().GetString("project")
			if project == "" {
				return errors.New("project should not be empty")
			}

			ids := strhlp.Map(strings.TrimSpace, args)
			if strhlp.Search("", ids) != -1 {
				return errors.New("task id/name should not be empty")
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

				if ids, err = search.GetTasksByName(
					c,
					api.GetTasksParam{
						Workspace: workspace,
						ProjectID: project,
						Active:    true,
					},
					ids,
				); err != nil {
					var errNF search.ErrNotFound
					if errors.As(err, &errNF) {
						return errors.New(
							"No active task with id or name containing '" +
								errNF.Reference + "' was found")
					}
					return err
				}
			}

			tasks := make([]dto.Task, len(ids))
			var g errgroup.Group
			for i := 0; i < len(ids); i++ {
				j := i
				g.Go(func() error {
					t, err := c.GetTask(api.GetTaskParam{
						Workspace: workspace,
						ProjectID: project,
						TaskID:    ids[j],
					})
					if err != nil {
						return err
					}

					tasks[j], err = c.UpdateTask(api.UpdateTaskParam{
						Workspace: workspace,
						ProjectID: t.ProjectID,
						TaskID:    t.ID,
						Name:      t.Name,
						Status:    api.TaskStatusDone,
					})

					return err
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			return util.TaskReport(cmd, of, tasks...)
		},
	}

	cmdutil.AddProjectFlags(cmd, f)
	util.TaskAddReportFlags(cmd, &of)
	return cmd
}
