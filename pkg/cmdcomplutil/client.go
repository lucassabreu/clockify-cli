package cmdcomplutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

// NewClientAutoComplete will provide auto-completion to flags or args
func NewClientAutoComplete(f factory) cmdcompl.SuggestFn {
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
		clients, err := c.GetClients(api.GetClientsParam{
			Workspace:       w,
			Archived:        &b,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for _, client := range clients {
			if toComplete != "" && !strings.Contains(client.ID, toComplete) {
				continue
			}
			va.Set(client.ID, client.Name)
		}

		return va, nil
	}
}
