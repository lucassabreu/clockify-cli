package cmdcomplutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

// NewProjectAutoComplete will provide auto-completion to flags or args
func NewProjectAutoComplete(f factory) cmdcompl.SuggestFn {
	return func(
		cmd *cobra.Command, args []string, toComplete string,
	) (cmdcompl.ValidArgs, error) {
		w, err := f.GetWorkspaceID()
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		c, err := f.Client()
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		b := false
		projects, err := c.GetProjects(api.GetProjectsParam{
			Workspace:       w,
			Archived:        &b,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for _, e := range projects {
			if toComplete != "" && !strings.Contains(e.ID, toComplete) {
				continue
			}
			va.Set(e.ID, e.Name)
		}

		return va, nil
	}
}
