package timehlp

import "time"

// TruncateDate clears the hours, minutes and seconds of a time.Time
func TruncateDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Today will return a UTC time.Time for the same day as time.Now() in Local
// time, but at 0:00:00.000
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
