package util

import (
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// DescriptionSuggestFn provides suggestions to user when setting the description of a
// time entry
type DescriptionSuggestFn func(string) []string

// descriptionCompleter looks for similar descriptions for auto-compliance
type descriptionCompleter struct {
	client       api.Client
	loaded       bool
	param        api.GetUserTimeEntriesParam
	descriptions []string
}

// NewDescriptionCompleter create or not a descriptionCompleter based on params
func NewDescriptionCompleter(f cmdutil.Factory) DescriptionSuggestFn {
	if !f.Config().GetBool(cmdutil.CONF_DESCR_AUTOCOMP) {
		return func(s string) []string { return []string{} }
	}

	workspaceID, err := f.GetWorkspaceID()
	if err != nil {
		return func(s string) []string { return []string{} }
	}

	userID, err := f.GetUserID()
	if err != nil {
		return func(s string) []string { return []string{} }
	}

	c, err := f.Client()
	if err != nil {
		return func(s string) []string { return []string{} }
	}

	end := time.Now().UTC()
	start := end.Add(time.Hour *
		time.Duration(-24*f.Config().GetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS)))

	d := &descriptionCompleter{
		client: c,
		param: api.GetUserTimeEntriesParam{
			Workspace: workspaceID,
			UserID:    userID,
			End:       &end,
			Start:     &start,
		},
	}

	return d.suggestFn
}

// getDescriptions load descriptions from recent time entries and list than
// unique ones
func (dc *descriptionCompleter) getDescriptions() []string {
	if dc.loaded {
		return dc.descriptions
	}

	tes, err := dc.client.GetUserTimeEntries(dc.param)

	dc.loaded = true
	if err != nil {
		return dc.descriptions
	}

	ss := []string{}
	for _, t := range tes {
		ss = append(ss, t.Description)
	}

	dc.descriptions = strhlp.Unique(ss)
	return dc.descriptions
}

// suggestFn returns a list of suggested descriptions based on a input string
func (dc *descriptionCompleter) suggestFn(toComplete string) []string {
	toComplete = strings.TrimSpace(toComplete)
	if toComplete == "" {
		return dc.getDescriptions()
	}

	toComplete = strhlp.Normalize(toComplete)
	return strhlp.Filter(
		func(s string) bool {
			return strings.Contains(strhlp.Normalize(s), toComplete)
		},
		dc.getDescriptions(),
	)
}

func newDescriptionAutoComplete(f cmdutil.Factory) cmdcompl.SuggestFn {
	return func(
		_ *cobra.Command, _ []string, toComplete string,
	) (cmdcompl.ValidArgs, error) {
		if !f.Config().GetBool(cmdutil.CONF_DESCR_AUTOCOMP) {
			return cmdcompl.EmptyValidArgs(), nil
		}

		dc := NewDescriptionCompleter(f)
		return cmdcompl.ValidArgsSlide(dc(toComplete)), nil
	}
}
