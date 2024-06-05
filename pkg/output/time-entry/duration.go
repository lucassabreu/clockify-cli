package timeentry

import (
	"fmt"
	"io"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func timeEntriesTotalDurationOnly(
	f func(time.Duration) string,
	timeEntries []dto.TimeEntry,
	w io.Writer,
) error {
	_, err := fmt.Fprintln(w, f(sumTimeEntriesDuration(timeEntries)))
	return err
}

// TimeEntriesTotalDurationOnlyAsFloat will only print the total duration as
// float
func TimeEntriesTotalDurationOnlyAsFloat(
	timeEntries []dto.TimeEntry, w io.Writer,
	l language.Tag) error {
	p := message.NewPrinter(l)
	println(l.String())

	return timeEntriesTotalDurationOnly(
		func(d time.Duration) string {
			return p.Sprintf("%f", number.Decimal(d.Hours()))
		},
		timeEntries,
		w,
	)
}

// TimeEntriesTotalDurationOnlyFormatted will only print the total duration as
// float
func TimeEntriesTotalDurationOnlyFormatted(
	timeEntries []dto.TimeEntry, w io.Writer) error {
	return timeEntriesTotalDurationOnly(
		durationToString,
		timeEntries,
		w,
	)
}
