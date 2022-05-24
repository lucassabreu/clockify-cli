package timehlp

import "time"

// TruncateDate resets the hours, minutes and seconds of a time.Time
func TruncateDate(t time.Time) time.Time {
	return t.Truncate(time.Hour * 24)
}
