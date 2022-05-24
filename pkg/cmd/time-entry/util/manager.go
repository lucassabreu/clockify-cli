package util

import "github.com/lucassabreu/clockify-cli/api/dto"

type CallbackFn func(dto.TimeEntryImpl) (dto.TimeEntryImpl, error)

func nullCallback(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
	return te, nil
}

// ManageEntry will runs all callback functions over the time entry, keeping
// the changes and returning it after
func ManageEntry(te dto.TimeEntryImpl, cbs ...CallbackFn) (
	dto.TimeEntryImpl, error) {
	return composeCallbacks(cbs...)(te)
}

func composeCallbacks(cbs ...CallbackFn) CallbackFn {
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
