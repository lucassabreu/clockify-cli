package timehlp

import "time"

// TruncateDate clears the hours, minutes and seconds of a time.Time for UTC
func TruncateDate(t time.Time) time.Time {
	return TruncateDateWithTimezone(t, time.UTC)
}

// TruncateDateWithTimezone clears the hours, minutes and seconds of a
// time.Time for a time.Location
func TruncateDateWithTimezone(t time.Time, l *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, l).
		Truncate(time.Second)
}

// Today will return a UTC time.Time for the same day as time.Now() in Local
// time, but at 0:00:00.000
func Today() time.Time {
	n := Now()
	return TruncateDateWithTimezone(n, n.Location())
}

// Now returns a time.Time using the local timezone
func Now() time.Time {
	return time.Now().In(time.Local).Truncate(time.Second)
}
