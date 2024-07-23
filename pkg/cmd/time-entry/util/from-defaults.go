package util

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

// FromDefaults starts a TimeEntryDTO with the current defaults
func FromDefaults(f cmdutil.Factory) Step {
	return func(ted TimeEntryDTO) (TimeEntryDTO, error) {
		d, err := f.TimeEntryDefaults().Read()
		if err != nil && err != defaults.DefaultsFileNotFoundErr {
			return ted, err
		}

		ted.ProjectID = d.ProjectID
		ted.TaskID = d.TaskID
		ted.TagIDs = d.TagIDs
		ted.Billable = d.Billable

		return ted, nil
	}
}
