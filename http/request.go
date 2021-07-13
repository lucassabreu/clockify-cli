package http

import (
	"strconv"
	"time"
)

func NewDateTime(t time.Time) DateTime {
	return DateTime{Time: t}
}

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
