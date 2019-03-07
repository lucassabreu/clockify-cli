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
	return []byte(strconv.Quote(d.Format("2006-01-02T15:04:05Z"))), nil
}

// TimeEntryStartEndRequest to get entries by range
type TimeEntryStartEndRequest struct {
	Start DateTime `json:"start"`
	End   DateTime `json:"end"`
}
