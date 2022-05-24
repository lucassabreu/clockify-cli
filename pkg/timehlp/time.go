package timehlp

import (
	"fmt"
	"strings"
	"time"
)

const (
	FullTimeFormat        = "2006-01-02 15:04:05"
	SimplerTimeFormat     = "2006-01-02 15:04"
	OnlyTimeFormat        = "15:04:05"
	SimplerOnlyTimeFormat = "15:04"
	NowTimeFormat         = "now"
)

func ConvertToTime(timeString string) (t time.Time, err error) {
	timeString = strings.ToLower(strings.TrimSpace(timeString))

	if NowTimeFormat == timeString {
		return time.Now().In(time.Local), nil
	}

	if strings.HasPrefix(timeString, "+") ||
		strings.HasPrefix(timeString, "-") {
		return relativeToTime(timeString)
	}

	if strings.HasPrefix(timeString, "yesterday ") {
		timeString = time.Now().Format("2006-01-02") + " " + timeString[10:]
	}

	l := len(timeString)
	if len(FullTimeFormat) != l &&
		len(SimplerTimeFormat) != l &&
		len(OnlyTimeFormat) != l &&
		len(SimplerOnlyTimeFormat) != l {
		return t, fmt.Errorf(
			"supported formats are: %s",
			strings.Join(
				[]string{
					FullTimeFormat, SimplerTimeFormat, OnlyTimeFormat,
					SimplerOnlyTimeFormat, NowTimeFormat,
				},
				", ",
			),
		)
	}

	if len(SimplerOnlyTimeFormat) == l || len(SimplerTimeFormat) == l {
		timeString = timeString + ":00"
		l = l + 3
	}

	if len(OnlyTimeFormat) == l {
		timeString = time.Now().Format("2006-01-02") + " " + timeString
	}

	return time.ParseInLocation(FullTimeFormat, timeString, time.Local)
}
