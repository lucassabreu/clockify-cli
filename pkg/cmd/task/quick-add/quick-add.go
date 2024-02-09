package quickadd

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// NewCmdQuickAdd will add multiple tasks to a project, but only setting its
// name
func NewCmdQuickAdd(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, []dto.Task) error,
) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:     "quick-add <name>...",
		Aliases: []string{"quick"},
		Short:   "Adds tasks to a project on Clockify",
		Args:    cmdutil.RequiredNamedArgs("name"),
		Long: "Adds a new active tasks to a project on Clockify, " +
			"but only allow setting their names.",
		Example: heredoc.Docf(`
			$ %[1]s -p special "Very Important"
			+--------------------------+----------------+--------+
			|            ID            |      NAME      | STATUS |
			+--------------------------+----------------+--------+
			| 62aa5d7049445270d7b979d6 | Very Important | ACTIVE |
			+--------------------------+----------------+--------+

			$ %[1]s -p special "Very Cool" -q
			dddddddddddddddddddddddd

			$ %[1]s -p special Billable "Not Billable" --csv
			id,name,status
			62ab145ec22de9759e6f6e35,Billable,ACTIVE
			62ab145ec22de9759e6f6e36,Not Billable,ACTIVE
		`, "clockify-cli task quick-add"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			p, _ := cmd.Flags().GetString("project")
			if f.Config().IsAllowNameForID() {
				if p, err = search.GetProjectByName(c, w, p, ""); err != nil {
					return err
				}
			}

			names := strhlp.Unique(args)
			tasks := make([]dto.Task, len(names))
			g := errgroup.Group{}
			for j := range names {
				i := j
				g.Go(func() error {
					tasks[i], err = c.AddTask(api.AddTaskParam{
						Workspace: w,
						ProjectID: p,
						Name:      names[i],
					})
					return err
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			if report != nil {
				return report(cmd.OutOrStdout(), &of, tasks)
			}

			return util.TaskReport(cmd, of, tasks...)
		},
	}

	cmdutil.AddProjectFlags(cmd, f)
	util.TaskAddReportFlags(cmd, &of)
	return cmd
}
