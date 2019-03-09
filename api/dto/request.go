package dto

import (
	"strconv"
	"time"
)

// DateTime is a time presentation for parameters
type DateTime struct {
	time.Time
}

// MarshalJSON converts DateTime correctly
func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.Time.UTC().Format("2006-01-02T15:04:05Z"))), nil
}

// TimeEntryStartEndRequest to get entries by range
type TimeEntryStartEndRequest struct {
	Start DateTime `json:"start"`
	End   DateTime `json:"end"`
}

// OutTimeEntryRequest to end the current time entry
type OutTimeEntryRequest struct {
	End DateTime `json:"end"`
}

// CreateTimeEntryRequest to create a time entry is created
type CreateTimeEntryRequest struct {
	Start       DateTime  `json:"start,omitempty"`
	End         *DateTime `json:"end,omitempty"`
	Billable    bool      `json:"billable,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectID   string    `json:"projectId,omitempty"`
	TaskID      string    `json:"taskId,omitempty"`
	TagIDs      []string  `json:"tagIds,omitempty"`
}

// UpdateTimeEntryRequest to update a time entry
type UpdateTimeEntryRequest struct {
	Start       DateTime  `json:"start,omitempty"`
	End         *DateTime `json:"end,omitempty"`
	Billable    bool      `json:"billable,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectID   string    `json:"projectId,omitempty"`
	TaskID      string    `json:"taskId,omitempty"`
	TagIDs      []string  `json:"tagIds,omitempty"`
}
