package timehlp

import "time"

// GetMonthRange given a time it returns the first and last date of a month
func GetMonthRange(ref time.Time) (first, last time.Time) {
	first = ref.AddDate(0, 0, ref.Day()*-1+1)
	last = first.AddDate(0, 1, -1)

	return
}

// GetWeekRange given a time it returns the first and last date of a week
func GetWeekRange(ref time.Time) (first, last time.Time) {
	first = ref.AddDate(0, 0, int(ref.Weekday())*-1)
	last = first.AddDate(0, 0, 7)

	return
}
