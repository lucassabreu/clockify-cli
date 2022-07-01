package timeentry

import (
	"fmt"
	"io"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
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
func TimeEntriesTotalDurationOnlyAsFloat(timeEntries []dto.TimeEntry, w io.Writer) error {
	return timeEntriesTotalDurationOnly(
		func(d time.Duration) string { return fmt.Sprintf("%f", d.Hours()) },
		timeEntries,
		w,
	)
}

// TimeEntryTotalDurationOnlyFormatted will only print the total duration as
// float
func TimeEntriesTotalDurationOnlyFormatted(
	timeEntries []dto.TimeEntry, w io.Writer) error {
	return timeEntriesTotalDurationOnly(
		durationToString,
		timeEntries,
		w,
	)
}
