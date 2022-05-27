package cmdcomplutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

// NewWorspaceAutoComplete will provice auto-completion for flags or args
func NewWorspaceAutoComplete(f factory) cmdcompl.SuggestFn {
	return func(
		cmd *cobra.Command, args []string, toComplete string,
	) (cmdcompl.ValidArgs, error) {
		c, err := f.Client()
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		ws, err := c.GetWorkspaces(api.GetWorkspaces{})

		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for i := range ws {
			if toComplete != "" && !strings.Contains(ws[i].ID, toComplete) {
				continue
			}
			va.Set(ws[i].ID, ws[i].Name)
		}

		return va, nil
	}
}
