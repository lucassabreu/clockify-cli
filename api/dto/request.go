package dto

import (
	"net/url"
	"strconv"
	"time"
)

// DateTime is a time presentation for parameters
type DateTime struct {
	time.Time
}

// MarshalJSON converts DateTime correctly
func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

func (d DateTime) String() string {
	return d.Time.UTC().Format("2006-01-02T15:04:05Z")
}

// TimeEntryStartEndRequest to get entries by range
type TimeEntryStartEndRequest struct {
	Start DateTime
	End   DateTime
}

func (r TimeEntryStartEndRequest) AppendToQuery(u url.URL) url.URL {
	v := u.Query()
	v.Add("start", r.Start.String())
	v.Add("end", r.End.String())
	u.RawQuery = v.Encode()

	return u
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
