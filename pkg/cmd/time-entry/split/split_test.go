package split_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/split"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdSplitShouldFail(t *testing.T) {
	w := dto.Workspace{ID: "w"}
	start, _ := timehlp.ConvertToTime("08:15")
	end, _ := timehlp.ConvertToTime("12:15")
	te := dto.TimeEntryImpl{
		WorkspaceID: w.ID,
		ID:          "timeentryid",
		Description: "Something",
		ProjectID:   "oldproj",
		TaskID:      "oldtask",
		TimeInterval: dto.TimeInterval{
			Start: start,
			End:   &end,
		},
	}

	findTT := func(t *testing.T) cmdutil.Factory {
		f := mocks.NewMockFactory(t)

		f.EXPECT().GetWorkspaceID().Return(w.ID, nil)
		f.EXPECT().GetUserID().Return(w.ID, nil)

		c := mocks.NewMockClient(t)
		f.EXPECT().Client().Return(c, nil)

		c.EXPECT().GetTimeEntry(api.GetTimeEntryParam{
			Workspace:   te.WorkspaceID,
			TimeEntryID: te.ID,
		}).
			Return(&te, nil)

		return f
	}

	tts := []struct {
		name string
		args []string
		f    func(*testing.T) cmdutil.Factory
		err  string
	}{
		{
			name: "time string is not valid",
			f: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
			args: []string{te.ID, "ff"},
			err:  "argument 2 could not be converted",
		},
		{
			name: "third time string is not valid",
			f: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
			args: []string{te.ID, "11:00", "ff"},
			err:  "argument 3 could not be converted",
		},
		{
			name: "split must be in order",
			f: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
			args: []string{te.ID, "8:30", "08:20"},
			err:  "splits must be in increasing order",
		},
		{
			name: "time entry not found",
			f: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				f.EXPECT().GetWorkspaceID().Return(w.ID, nil)
				f.EXPECT().GetUserID().Return(w.ID, nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetTimeEntry(api.GetTimeEntryParam{
					Workspace:   te.WorkspaceID,
					TimeEntryID: te.ID,
				}).
					Return(nil, nil)

				return f
			},
			args: []string{te.ID, "11:00"},
			err:  "not found",
		},
		{
			name: "split must be after the time entry start",
			f:    findTT,
			args: []string{te.ID, "7:00"},
			err:  "time splits must be after .* 08:15",
		},
		{
			name: "split must be before the time entry ends",
			f:    findTT,
			args: []string{te.ID, "18:00"},
			err:  "time splits must be before .* 12:15",
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			called := false
			cmd := split.NewCmdSplit(tt.f(t), func(
				_ []dto.TimeEntry, _ io.Writer, _ util.OutputFlags) error {
				called = true
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			assert.False(t, called)
			assert.Error(t, err)
			assert.Regexp(t, tt.err, err.Error())
		})
	}
}

func TestNewCmdSplitShouldCreateNewEntries(t *testing.T) {
	w := dto.Workspace{ID: "w"}
	start, _ := timehlp.ConvertToTime("08:00")
	te := dto.TimeEntryImpl{
		WorkspaceID: w.ID,
		ID:          "timeentryid",
		Description: "Something",
		ProjectID:   "oldproj",
		TaskID:      "oldtask",
		TimeInterval: dto.TimeInterval{
			Start: start,
		},
	}

	getH := func(c *mocks.MockClient, id string) {

		c.EXPECT().GetHydratedTimeEntry(api.GetTimeEntryParam{
			Workspace:   w.ID,
			TimeEntryID: id,
		}).
			Return(&dto.TimeEntry{
				WorkspaceID: te.WorkspaceID,
				ID:          id,
			}, nil)
	}

	getAndUpdate := func(
		t *testing.T, te dto.TimeEntryImpl, end time.Time,
	) (cmdutil.Factory, *mocks.MockClient) {
		f := mocks.NewMockFactory(t)

		f.EXPECT().GetWorkspaceID().Return(w.ID, nil)
		f.EXPECT().GetUserID().Return(w.ID, nil)

		c := mocks.NewMockClient(t)
		f.EXPECT().Client().Return(c, nil)

		c.EXPECT().GetTimeEntry(api.GetTimeEntryParam{
			Workspace:   te.WorkspaceID,
			TimeEntryID: te.ID,
		}).
			Return(&te, nil)

		c.EXPECT().UpdateTimeEntry(api.UpdateTimeEntryParam{
			Workspace:   te.WorkspaceID,
			TimeEntryID: te.ID,
			Start:       te.TimeInterval.Start,
			End:         &end,
			Billable:    false,
			Description: te.Description,
			ProjectID:   te.ProjectID,
			TaskID:      te.TaskID,
			TagIDs:      te.TagIDs,
		}).
			Return(te, nil)

		getH(c, te.ID)

		return f, c
	}

	create := func(c *mocks.MockClient, id string, ted util.TimeEntryDTO) {
		c.EXPECT().CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   te.WorkspaceID,
			Start:       ted.Start,
			End:         ted.End,
			Description: te.Description,
			ProjectID:   te.ProjectID,
			TaskID:      te.TaskID,
			TagIDs:      te.TagIDs,
			Billable:    &te.Billable,
		}).
			Return(dto.TimeEntryImpl{ID: id}, nil)

		getH(c, id)
	}

	tts := []struct {
		name string
		args []string
		f    func(*testing.T) cmdutil.Factory
	}{
		{
			name: "split in two",
			args: []string{te.ID, "8:30"},
			f: func(t *testing.T) cmdutil.Factory {
				s, _ := timehlp.ConvertToTime("08:30")
				f, c := getAndUpdate(t, te, s)

				create(c, "123", util.TimeEntryDTO{
					Workspace:   te.ID,
					UserID:      te.UserID,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					Description: te.Description,
					Start:       s,
					End:         nil,
				})

				return f
			},
		},
		{
			name: "split in three",
			args: []string{te.ID, "08:20", "8:30"},
			f: func(t *testing.T) cmdutil.Factory {
				s, _ := timehlp.ConvertToTime("08:20")
				f, c := getAndUpdate(t, te, s)

				e, _ := timehlp.ConvertToTime("08:30")
				create(c, "123", util.TimeEntryDTO{
					Workspace:   te.ID,
					UserID:      te.UserID,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					Description: te.Description,
					Start:       s,
					End:         &e,
				})

				create(c, "456", util.TimeEntryDTO{
					Workspace:   te.ID,
					UserID:      te.UserID,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					Description: te.Description,
					Start:       e,
					End:         nil,
				})

				return f
			},
		},
		{
			name: "split in three with end",
			args: []string{te.ID, "08:30", "9:00"},
			f: func(t *testing.T) cmdutil.Factory {
				end, _ := timehlp.ConvertToTime("10:00")
				te := dto.TimeEntryImpl{
					WorkspaceID: w.ID,
					ID:          "timeentryid",
					Description: "Something",
					ProjectID:   "oldproj",
					TaskID:      "oldtask",
					TimeInterval: dto.NewTimeInterval(
						te.TimeInterval.Start,
						&end,
					),
				}

				s, _ := timehlp.ConvertToTime("08:30")
				f, c := getAndUpdate(t, te, s)

				e, _ := timehlp.ConvertToTime("09:00")
				create(c, "123", util.TimeEntryDTO{
					Workspace:   te.ID,
					UserID:      te.UserID,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					Description: te.Description,
					Start:       s,
					End:         &e,
				})

				create(c, "456", util.TimeEntryDTO{
					Workspace:   te.ID,
					UserID:      te.UserID,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					Description: te.Description,
					Start:       e,
					End:         &end,
				})

				return f
			},
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			called := false
			cmd := split.NewCmdSplit(tt.f(t), func(
				_ []dto.TimeEntry, _ io.Writer, _ util.OutputFlags) error {
				called = true
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if assert.NoError(t, err) {
				return
			}

			if assert.True(t, called) {
				return
			}
		})
	}
}
