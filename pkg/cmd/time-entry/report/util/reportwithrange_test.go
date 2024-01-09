package util_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func newDate(s string) time.Time {
	date, _ := time.ParseInLocation("2006-01-02", s, time.UTC)
	return date
}

func TestReportWithRange(t *testing.T) {
	date := newDate("2006-01-02")
	first := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	last := first.AddDate(0, 0, 3)
	tts := []struct {
		name     string
		factory  func(*testing.T) cmdutil.Factory
		flags    func(*testing.T) util.ReportFlags
		expected string
		err      string
	}{
		{
			name: "no user",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("", errors.New("no user"))
				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				return util.NewReportFlags()
			},
			err: "no user",
		},
		{
			name: "no workspace",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("", errors.New("no workspace"))
				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				return util.NewReportFlags()
			},
			err: "no workspace",
		},
		{
			name: "no client",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)
				f.On("Client").Return(nil, errors.New("no client"))
				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				return util.NewReportFlags()
			},
			err: "no client",
		},
		{
			name: "http error project",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{}, errors.New("http error"))

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Project = "p"
				return rf
			},
			err: "http error",
		},
		{
			name: "invalid project",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{{Name: "right"}}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Project = "wrong"
				return rf
			},
			err: "No project.*wrong' was found",
		},
		{
			name: "range http error",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{{ID: "p", Name: "right"}}, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.TimeEntry{}, errors.New("http error"))

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Project = "right"
				return rf
			},
			err: "http error",
		},
		{
			name: "project and description",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(false)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p",
					Description:     "desc",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "time-entry-1",
						TimeInterval: dto.TimeInterval{Start: last}},
					{ID: "time-entry-2",
						TimeInterval: dto.TimeInterval{Start: first}},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Project = "p"
				rf.Description = "desc"
				rf.Quiet = true
				return rf
			},
			expected: heredoc.Doc(`
				time-entry-2
				time-entry-1
			`),
		},
		{
			name: "fill missing dates",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "time-entry-1", TimeInterval: dto.TimeInterval{
						Start: newDate("2006-01-04")}},
					{ID: "time-entry-2", TimeInterval: dto.TimeInterval{
						Start: newDate("2006-01-01")}},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.TimeZone = "UTC"
				rf.FillMissingDates = true
				rf.Format = "{{.ID}};{{ .TimeInterval.Start.Format " +
					`"2006-01-02"` +
					" }}"
				return rf
			},
			expected: heredoc.Doc(`
				time-entry-2;2006-01-01
				;2006-01-02
				;2006-01-03
				time-entry-1;2006-01-04
			`),
		},
		{
			name: "billable only",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "time-entry-1", Billable: true},
					{ID: "time-entry-2", Billable: false},
					{ID: "time-entry-3", Billable: true},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Billable = true
				rf.Quiet = true
				return rf
			},
			expected: heredoc.Doc(`
				time-entry-1
				time-entry-3
			`),
		},
		{
			name: "not billable only",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "time-entry-1", Billable: true},
					{ID: "time-entry-2", Billable: false},
					{ID: "time-entry-3", Billable: true},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.NotBillable = true
				rf.Quiet = true
				return rf
			},
			expected: heredoc.Doc(`
				time-entry-2
			`),
		},
		{
			name: "not billable & tag cli only",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				tag := dto.Tag{ID: "t1", Name: "Client"}
				c.On("GetTags", api.GetTagsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Tag{tag}, nil)

				c.On("LogRange", api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					FirstDate:       first,
					LastDate:        last,
					TagIDs:          []string{tag.ID},
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "te-1", Tags: []dto.Tag{tag}, Billable: true},
					{ID: "te-2", Tags: []dto.Tag{tag}, Billable: false},
					{ID: "te-3", Tags: []dto.Tag{tag}, Billable: true},
					{ID: "te-4", Tags: []dto.Tag{tag}, Billable: false},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.NotBillable = true
				rf.Quiet = true
				rf.TagIDs = []string{"cli"}
				return rf
			},
			expected: heredoc.Doc(`
				te-2
				te-4
			`),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			err := util.ReportWithRange(
				tt.factory(t),
				date,
				date.AddDate(0, 0, 2),
				b,
				tt.flags(t),
			)

			if tt.err != "" {
				if assert.Error(t, err) {
					assert.Regexp(t, tt.err, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, b.String())
		})
	}
}
