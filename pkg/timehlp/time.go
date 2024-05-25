package timehlp

import (
	"fmt"
	"strings"
	"time"
)

const (
	FullTimeFormat           = "2006-01-02 15:04:05"
	SimplerTimeFormat        = "2006-01-02 15:04"
	OnlyTimeFormat           = "15:04:05"
	SimplerOnlyTimeFormat    = "15:04"
	SimplerOnlyTimeFormatWL  = "5:04"
	NowTimeFormat            = "now"
	SimplestOnlyTimeFormat   = "1504"
	SimplestOnlyTimeFormatWL = "504"
)

// ConvertToTime will try to convert a string do time.Time looking for the
// format that best fits it and assuming "today" when necessary.
// If the string starts with `yesterday`, than it will be exchanged for a
// date-string with the format: 2006-01-02
// If the string starts with `+` or `-` than the string will be treated as
// "relative time expressions", and will be calculated as the diff from now and
// it.
// If the string is "now" than `time.Now()` in the local timezone will be
// returned.
func ConvertToTime(timeString string) (t time.Time, err error) {
	timeString = strings.ToLower(strings.TrimSpace(timeString))

	if NowTimeFormat == timeString {
		return Now(), nil
	}

	if strings.HasPrefix(timeString, "+") ||
		strings.HasPrefix(timeString, "-") {
		return relativeToTime(timeString)
	}

	if strings.HasPrefix(timeString, "yesterday ") {
		timeString = Today().
			Add(-1).Format("2006-01-02") + " " + timeString[10:]
	}

	l := len(timeString)
	if len(FullTimeFormat) != l &&
		len(SimplerTimeFormat) != l &&
		len(OnlyTimeFormat) != l &&
		len(SimplerOnlyTimeFormat) != l &&
		len(SimplestOnlyTimeFormat) != l &&
		len(SimplestOnlyTimeFormatWL) != l {
		return t, fmt.Errorf(
			"supported formats are: %s",
			strings.Join(
				[]string{
					FullTimeFormat, SimplerTimeFormat, OnlyTimeFormat,
					SimplerOnlyTimeFormat, SimplerOnlyTimeFormatWL, NowTimeFormat,
					SimplestOnlyTimeFormat, SimplestOnlyTimeFormatWL,
				},
				", ",
			),
		)
	}

	timeString = normalizeFormats(timeString)
	t, err = time.ParseInLocation(FullTimeFormat, timeString, time.Local)
	if err != nil {
		return t, err
	}

	return t.Truncate(time.Second), nil
}

// Adds data to the partial timeString to match a full
// datetime with seconds precission.
// Receives a time in any of the defined formats, and return
// a date in the FullTimeFormat
func normalizeFormats(timeString string) string {
	l := len(timeString)

	// change from 9:14 to 09:14
	if len(SimplerOnlyTimeFormatWL) == l && strings.Contains(timeString, ":") {
		timeString = "0" + timeString
		l = l + 1
	}

	// change from 914 to 0914
	if len(SimplestOnlyTimeFormatWL) == l && !strings.Contains(timeString, ":") {
		timeString = "0" + timeString
		l = l + 1
	}

	// change from 0914 to 09:14
	if len(SimplestOnlyTimeFormat) == l {
		timeString = timeString[0:2] + ":" + timeString[2:]
		l = l + 1
	}

	// change from 09:14 to 09:14:00
	if len(SimplerOnlyTimeFormat) == l || len(SimplerTimeFormat) == l {
		timeString = timeString + ":00"
		l = l + 3
	}

	// change from 09:14 to 2006-01-02 09:14:00
	if len(OnlyTimeFormat) == l {
		timeString = Today().Format("2006-01-02") + " " + timeString
	}
	return timeString
}
