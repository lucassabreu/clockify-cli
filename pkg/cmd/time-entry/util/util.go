// util package provides reusable functionality to the commands under
// pkg/cmd/time-entry, be it editing, creating, or rendering time entries
package util

import (
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TimeEntryDTO is used to keep and update the data of a time entry before
// changing it, taking into account optional values (nil)
type TimeEntryDTO struct {
	ID          string
	Workspace   string
	UserID      string
	ProjectID   string
	Client      string
	TaskID      string
	Description string
	Start       time.Time
	End         *time.Time
	TagIDs      []string
	Billable    *bool
	Locked      *bool
}

// Step is used to stack multiple actions to be executed over a TimeEntryDTO
type Step func(TimeEntryDTO) (TimeEntryDTO, error)

func skip(te TimeEntryDTO) (TimeEntryDTO, error) {
	return te, nil
}

// Do will runs all callback functions over the time entry, keeping
// the changes and returning it after
func Do(te TimeEntryDTO, cbs ...Step) (TimeEntryDTO, error) {
	return compose(cbs...)(te)
}

func compose(cbs ...Step) Step {
	return func(dto TimeEntryDTO) (TimeEntryDTO, error) {
		var err error
		for _, cb := range cbs {
			if dto, err = cb(dto); err != nil {
				return dto, err
			}
		}

		return dto, err
	}
}

// TimeEntryImplToDTO returns a TimeEntryDTO using the information from a
// TimeEntryImpl
func TimeEntryImplToDTO(t dto.TimeEntryImpl) TimeEntryDTO {
	return TimeEntryDTO{
		Workspace:   t.WorkspaceID,
		UserID:      t.UserID,
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		TaskID:      t.TaskID,
		Description: t.Description,
		Start:       t.TimeInterval.Start,
		End:         t.TimeInterval.End,
		TagIDs:      t.TagIDs,
		Billable:    &t.Billable,
		Locked:      &t.IsLocked,
	}
}

// TimeEntryDTOToImpl returns a TimeEntryImpl using the information from a
// TimeEntryDTO
func TimeEntryDTOToImpl(t TimeEntryDTO) dto.TimeEntryImpl {
	if t.Billable == nil {
		b := false
		t.Billable = &b
	}

	if t.Locked == nil {
		b := false
		t.Locked = &b
	}

	return dto.TimeEntryImpl{
		WorkspaceID:  t.Workspace,
		UserID:       t.UserID,
		Description:  t.Description,
		ID:           t.ID,
		ProjectID:    t.ProjectID,
		TagIDs:       t.TagIDs,
		TaskID:       t.TaskID,
		TimeInterval: dto.NewTimeInterval(t.Start, t.End),
		Billable:     *t.Billable,
		IsLocked:     *t.Locked,
	}
}
