package cmdcomplutil

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewProjectAutoComplete will provide auto-completion to flags or args
func NewProjectAutoComplete(f factory, config config) cmdcompl.SuggestFn {
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
		ps, err := c.GetProjects(api.GetProjectsParam{
			Workspace:       w,
			Archived:        &b,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return cmdcompl.EmptyValidArgs(), err
		}

		filter := makeFilter(toComplete, config)

		psf := make([]dto.Project, 0)
		padding := 0
		for i := range ps {
			if !filter(ps[i]) {
				continue
			}

			if padding < len(ps[i].Name) {
				padding = len(ps[i].Name)
			}

			psf = append(psf, ps[i])
		}

		format := func(p dto.Project) string { return p.Name }
		if config.IsSearchProjectWithClientsName() {
			f := fmt.Sprintf("%%-%ds", padding)
			format = func(p dto.Project) string {
				client := "Without Client"
				if p.ClientID != "" {
					client = p.ClientID + " -- " + p.ClientName
				}
				return fmt.Sprintf(f, p.Name) + " | " + client
			}
		}

		va := make(cmdcompl.ValidArgsMap)
		for i := range psf {
			va.Set(psf[i].ID, format(psf[i]))
		}

		return va, nil
	}
}

func makeFilter(toComplete string, config config) func(dto.Project) bool {
	if toComplete == "" {
		return func(_ dto.Project) bool { return true }
	}

	if config.IsAllowNameForID() &&
		config.IsSearchProjectWithClientsName() {
		s := strhlp.IsSimilar(toComplete)
		return func(p dto.Project) bool {
			return strings.Contains(p.ID, toComplete) || s(p.Name) ||
				strings.Contains(p.ClientID, toComplete) || s(p.ClientName)
		}
	}

	if config.IsAllowNameForID() {
		s := strhlp.IsSimilar(toComplete)
		return func(p dto.Project) bool {
			return strings.Contains(p.ID, toComplete) || s(p.Name)
		}
	}

	return func(p dto.Project) bool {
		return strings.Contains(p.ID, toComplete)
	}
}
