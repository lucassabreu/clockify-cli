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

		workspaces, err := c.GetWorkspaces(api.GetWorkspaces{})

		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for _, w := range workspaces {
			if toComplete != "" && !strings.Contains(w.ID, toComplete) {
				continue
			}
			va.Set(w.ID, w.Name)
		}

		return va, nil
	}
}
