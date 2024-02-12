package today_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/today"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestCmdToday(t *testing.T) {
	first := time.Now()
	first = time.Date(
		first.Year(),
		first.Month(),
		first.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	last := first.AddDate(0, 0, 1)

	tts := []struct {
		name     string
		args     string
		err      error
		expected string
		factory  func(*testing.T) cmdutil.Factory
	}{
		{
			name: "error on multi format",
			args: "--format {} --json --csv -q --md " +
				"--duration-float --duration-formatted",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
			err: errors.New(
				"the following flags can't be used together: " +
					"`csv`, `duration-float`, `duration-formatted`, " +
					"`format`, `json`, `md` and `quiet`",
			),
		},
		{
			name: "all of them, but only ids",
			args: "-q",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("user-id", nil)
				f.On("GetWorkspaceID").Return("w-id", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w-id",
					UserID:          "user-id",
					FirstDate:       first,
					LastDate:        last,
					TagIDs:          []string{},
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.TimeEntry{{ID: "time-entry-id"}},
						nil,
					)

				return f
			},
			expected: "time-entry-id\n",
		},
		{
			name: "all of them, but fails",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("user-id", nil)
				f.On("GetWorkspaceID").Return("w-id", nil)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w-id",
					UserID:          "user-id",
					FirstDate:       first,
					LastDate:        last,
					TagIDs:          []string{},
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.TimeEntry{},
						errors.New("failed"),
					)

				return f
			},
			err: errors.New("failed"),
		},
		{
			name: "only project x, no results",
			args: "--project x",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("user-id", nil)
				f.On("GetWorkspaceID").Return("w-id", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				bFalse := false
				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w-id",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.Project{{
							ID:   "project-id",
							Name: "xpecial",
						}},
						nil,
					)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w-id",
					UserID:          "user-id",
					FirstDate:       first,
					LastDate:        last,
					ProjectID:       "project-id",
					TagIDs:          []string{},
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.TimeEntry{},
						errors.New("failed"),
					)

				return f
			},
			err: errors.New("failed"),
		},
		{
			name: "only with desc on description",
			args: "--description desc -q",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("user-id", nil)
				f.On("GetWorkspaceID").Return("w-id", nil)

				f.On("Config").Return(mocks.NewMockConfig(t))

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w-id",
					UserID:          "user-id",
					FirstDate:       first,
					LastDate:        last,
					Description:     "desc",
					TagIDs:          []string{},
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.TimeEntry{
							{ID: "time-entry-1"},
							{ID: "time-entry-2"},
						},
						nil,
					)

				return f
			},
			expected: heredoc.Doc(`
				time-entry-1
				time-entry-2
			`),
		},
	}

	for i := range tts {
		tt := tts[i]
		t.Run(tt.name, func(t *testing.T) {
			cmd := today.NewCmdToday(tt.factory(t))
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			cmd.SetArgs(strings.Split(tt.args, " "))

			out := bytes.NewBufferString("")

			cmd.SetOut(out)
			cmd.SetErr(out)

			_, err := cmd.ExecuteC()

			if tt.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err.Error())
				return
			}

			assert.Equal(t, tt.expected, out.String())
			assert.NoError(t, err)
		})
	}
}
