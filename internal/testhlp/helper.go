package testhlp

import "time"

// MustParseTime will parse a string as time.Time or panic
func MustParseTime(l, v string) time.Time {
	t, err := time.Parse(l, v)
	if err == nil {
		return t
	}
	panic(err)
}
