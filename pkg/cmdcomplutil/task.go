package cmdcomplutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

// NewTaskAutoComplete will provide auto-completion to flags or args
func NewTaskAutoComplete(f factory) cmdcompl.SuggestFn {
	return func(
		cmd *cobra.Command, args []string, toComplete string,
	) (cmdcompl.ValidArgs, error) {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		w, err := f.GetWorkspaceID()
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		c, err := f.Client()
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		ts, err := c.GetTasks(api.GetTasksParam{
			Workspace: w,
			ProjectID: project,
			Active:    true,
		})
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for i := range ts {
			if toComplete != "" && !strings.Contains(ts[i].ID, toComplete) {
				continue
			}
			va.Set(ts[i].ID, ts[i].Name)
		}

		return va, nil
	}
}
