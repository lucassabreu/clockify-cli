package timehlp

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var ErrInvalidReliveTime = errors.New(
	"supported relative time formats are: " +
		"+15:04:05, +15:04 or unit descriptive +1d15h4m5s, " +
		"+15h5s, 120m",
)

func relativeToTime(timeString string) (t time.Time, err error) {
	var d time.Duration
	timeString = strings.ReplaceAll(timeString, " ", "")

	if c := strings.Count(timeString, ":"); c > 0 {
		d, err = relativeColonTimeToDuration(timeString[1:])
	} else {
		d, err = relativeUnitDescriptiveTimeToDuration(timeString[1:])
	}

	if timeString[0] == '-' {
		d = d * -1
	}

	t = Now().Add(d)
	return
}

func relativeColonTimeToDuration(s string) (d time.Duration, err error) {
	parts := strings.Split(s, ":")
	c := len(parts)
	if c > 2 || c == 0 {
		return d, ErrInvalidReliveTime
	}

	u := time.Second
	for i := c - 1; i >= 0; i-- {
		p := strings.TrimPrefix(parts[i], "0")
		v, err := strconv.Atoi(p)
		if err != nil && p != "" {
			return d, ErrInvalidReliveTime
		}
		d = d + time.Duration(v)*u
		u = u * 60
	}

	return
}

func relativeUnitDescriptiveTimeToDuration(s string) (
	d time.Duration, err error) {
	var u time.Duration
	var i, j int
	for ; i < len(s); i++ {
		switch s[i] {
		case 'd':
			u = time.Hour * 24
		case 'h':
			u = time.Hour
		case 'm':
			u = time.Minute
		case 's':
			u = time.Second
		default:
			continue
		}

		v, err := strconv.Atoi(s[j:i])
		if err != nil {
			return d, ErrInvalidReliveTime
		}

		d = d + time.Duration(v)*u
		j = i + 1
	}

	if i != j {
		return d, ErrInvalidReliveTime
	}

	return d, nil
}
