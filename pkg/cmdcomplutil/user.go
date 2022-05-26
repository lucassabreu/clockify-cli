package cmdcomplutil

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

// NewUserAutoComplete will provice auto-completion for flags or args
func NewUserAutoComplete(f factory) cmdcompl.SuggestFn {
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

		users, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
			Workspace: w,
		})

		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		va := make(cmdcompl.ValidArgsMap)
		toComplete = strings.ToLower(toComplete)
		for _, user := range users {
			if toComplete != "" && !strings.Contains(user.ID, toComplete) {
				continue
			}
			va.Set(user.ID, fmt.Sprintf("%s (%s)", user.Name, user.Email))
		}

		return va, nil
	}
}
