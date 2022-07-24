package timeentryhlp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/pkg/errors"
)

const (
	AliasCurrent = "current"
	AliasLast    = "last"
	AliasLatest  = "latest"
)

// GetLatestEntryEntry will return the last time entry of a user, if it exists
func GetLatestEntryEntry(
	c api.Client, workspace, userID string) (dto.TimeEntryImpl, error) {
	return GetTimeEntry(c, workspace, userID, AliasLatest)
}

var ErrNoTimeEntry = errors.New("time entry was not found")

func mayNotFound(tei *dto.TimeEntryImpl, err error) (
	dto.TimeEntryImpl, error) {
	if err != nil {
		return dto.TimeEntryImpl{}, err
	}

	if tei == nil {
		return dto.TimeEntryImpl{}, ErrNoTimeEntry
	}

	return *tei, nil
}

// GetTimeEntry will look for the time entry of a user for the id or alias
// provided
func GetTimeEntry(
	c api.Client,
	workspace,
	userID,
	id string,
) (dto.TimeEntryImpl, error) {
	id = strings.TrimSpace(strings.ToLower(id))

	var onlyInProgress *bool
	switch id {
	case "^0", AliasCurrent:
		tei, err := mayNotFound(c.GetTimeEntryInProgress(
			api.GetTimeEntryInProgressParam{
				Workspace: workspace,
				UserID:    userID,
			}))
		if err == ErrNoTimeEntry {
			return tei, errors.Wrap(err, "looking for running time entry")
		}

		return tei, err
	case "^1", AliasLast:
		id = AliasLast
		b := false
		onlyInProgress = &b
	case AliasLatest:
		id = AliasLatest
		onlyInProgress = nil
	}

	if id != AliasLast && id != AliasLatest && !strings.HasPrefix(id, "^") {
		return mayNotFound(c.GetTimeEntry(api.GetTimeEntryParam{
			Workspace:   workspace,
			TimeEntryID: id,
		}))
	}

	page := 1
	if strings.HasPrefix(id, "^") {
		var err error
		if page, err = strconv.Atoi(id[1:]); err != nil {
			return dto.TimeEntryImpl{}, fmt.Errorf(
				`n on "^n" must be a unsigned integer, you sent: %s`,
				id[1:],
			)
		}
	}

	list, err := c.GetUserTimeEntries(api.GetUserTimeEntriesParam{
		Workspace:      workspace,
		UserID:         userID,
		OnlyInProgress: onlyInProgress,
		PaginationParam: api.PaginationParam{
			PageSize: 1,
			Page:     page,
		},
	})

	if err != nil {
		return dto.TimeEntryImpl{}, err
	}

	if len(list) == 0 {
		return dto.TimeEntryImpl{}, ErrNoTimeEntry
	}

	return list[0], err
}
