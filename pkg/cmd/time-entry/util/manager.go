package util

import "github.com/lucassabreu/clockify-cli/api/dto"

type DoFn func(dto.TimeEntryImpl) (dto.TimeEntryImpl, error)

func nullCallback(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
	return te, nil
}

// Do will runs all callback functions over the time entry, keeping
// the changes and returning it after
func Do(te dto.TimeEntryImpl, cbs ...DoFn) (
	dto.TimeEntryImpl, error) {
	return compose(cbs...)(te)
}

func compose(cbs ...DoFn) DoFn {
	return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		var err error
		for _, cb := range cbs {
			if tei, err = cb(tei); err != nil {
				return tei, err
			}
		}

		return tei, err
	}
}
