package cmd

import (
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/viper"
)

type suggestFn func(string) []string

// descriptionCompleter looks for similar descriptions for auto-compliance
type descriptionCompleter struct {
	client       *api.Client
	loaded       bool
	param        api.GetUserTimeEntriesParam
	descriptions []string
}

// newDescriptionCompleter create or not a descriptionCompleter based on params
func newDescriptionCompleter(
	c *api.Client,
	workspaceID,
	userID string,
) suggestFn {
	if !viper.GetBool(DESCR_AUTOCOMP) {
		return func(s string) []string { return []string{} }
	}

	end := time.Now().UTC()
	start := end.Add(time.Hour *
		time.Duration(-24*viper.GetInt(DESCR_AUTOCOMP_DAYS)))

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
		func(s string) bool { return strings.Contains(strhlp.Normalize(s), toComplete) },
		dc.getDescriptions(),
	)
}
