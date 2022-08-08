package add

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAdd represents the add command
func NewCmdAdd(
	f cmdutil.Factory,
	report func(io.Writer, *util.OutputFlags, dto.Task) error,
) *cobra.Command {
	of := util.OutputFlags{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a new task to a project on Clockify",
		Long: heredoc.Doc(`
			Adds a new active task to a project on Clockify, also allows to assign users to it at the same time

			Tasks will be created as billable or not depending on the project settings.
			If you set a estimate for the task, but the project is set as manual estimation, then it will have no effect on Clockify.
		`),
		Example: heredoc.Docf(`
			$ %[1]s -p special --name="Very Important"
			+--------------------------+----------------+--------+
			|            ID            |      NAME      | STATUS |
			+--------------------------+----------------+--------+
			| 62aa5d7049445270d7b979d6 | Very Important | ACTIVE |
			+--------------------------+----------------+--------+

			$ %[1]s -p special --name="Very Cool" --assign john@example.com | \
			  jq '.[] |.assigneeIds' --compact-output
			["dddddddddddddddddddddddd"]

			$ %[1]s -p special --name Billable --billable --quiet
			62ab129e4ebb4f143c8e8622

			$ %[1]s -p special --name "Not Billable" --not-billable --csv
			id,name,status
			62ab145ec22de9759e6f6e36,Not Billable,ACTIVE

			$ %[1]s -p special --name 'With 1H to Make' --estimate 1
			+--------------------------+-----------------+--------+
			|            ID            |       NAME      | STATUS |
			+--------------------------+-----------------+--------+
			| 62aa5d7049445270d7b979d6 | With 1H to Make | ACTIVE |
			+--------------------------+-----------------+--------+
		`, "clockify-cli task add"),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			task, err := c.AddTask(api.AddTaskParam{
				Workspace:   fl.Workspace,
				ProjectID:   fl.ProjectID,
				Name:        fl.Name,
				Estimate:    fl.Estimate,
				AssigneeIDs: fl.AssigneeIDs,
				Billable:    fl.Billable,
			})
			if err != nil {
				return err
			}

			if report != nil {
				return report(cmd.OutOrStdout(), &of, task)
			}

			return util.TaskReport(cmd, of, task)
		},
	}

	util.TaskAddReportFlags(cmd, &of)

	util.TaskAddPropFlags(cmd, f)
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
