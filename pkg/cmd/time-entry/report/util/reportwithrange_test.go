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
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
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
				rf.Projects = []string{"p"}
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
				cf.On("IsSearchProjectWithClientsName").Return(false)

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
				rf.Projects = []string{"wrong"}
				return rf
			},
			err: "No project.*wrong' was found",
		},
		{
			name: "invalid client",
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
				rf.Client = "right"
				rf.Projects = []string{"wrong"}
				return rf
			},
			err: "No client.*right' was found",
		},
		{
			name: "invalid project for client",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)
				cf.On("IsSearchProjectWithClientsName").Return(false)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return(
						[]dto.Project{{
							Name:       "right",
							ClientName: "right",
							ClientID:   "r1",
						}},
						nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Client = "right"
				rf.Projects = []string{"wrong"}
				return rf
			},
			err: "No project.*wrong' was found for client 'right'",
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
				cf.On("IsSearchProjectWithClientsName").Return(false)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{
						{
							ID:         "p",
							Name:       "right",
							ClientName: "right",
							ClientID:   "c1",
						},
						{
							ID:         "p",
							Name:       "right",
							ClientName: "wrong",
							ClientID:   "c2",
						},
					}, nil)

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
				rf.Projects = []string{"right"}
				rf.Client = "right"
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

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
				})

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
				rf.Projects = []string{"p"}
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

				f.On("Config").Return(&mocks.SimpleConfig{})

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

				f.EXPECT().Config().Return(&mocks.SimpleConfig{})

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

				f.EXPECT().Config().Return(&mocks.SimpleConfig{})

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

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

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
		{
			name: "multiple projects",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				f.EXPECT().Config().Return(
					&mocks.SimpleConfig{AllowNameForID: true})

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p1", Name: "p1"},
					{ID: "p2", Name: "p2"},
					{ID: "p3", Name: "p3"},
				}, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p1",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "te-1",
						TimeInterval: dto.TimeInterval{
							Start: first,
						},
					},
					{ID: "te-3",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(2)),
						},
					},
				}, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p2",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "te-2",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(1)),
						},
					},
					{ID: "te-4",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(3)),
						},
					},
				}, nil)

				return f
			},
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Quiet = true
				rf.Projects = []string{"p1", "p2"}
				return rf
			},
			expected: heredoc.Doc(`
				te-1
				te-2
				te-3
				te-4
			`),
		},
		{
			name: "projects form a client",
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Quiet = true
				rf.Client = "me"
				return rf
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				f.EXPECT().Config().Return(
					&mocks.SimpleConfig{AllowNameForID: true})

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.EXPECT().GetClients(api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{
						{ID: "c1", Name: "me"},
						{ID: "c2", Name: "you"},
					}, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Clients:         []string{"c1"},
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p1", Name: "p1", ClientID: "c1", ClientName: "me"},
					{ID: "p3", Name: "p3", ClientID: "c1", ClientName: "me"},
				}, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p1",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "te-1",
						TimeInterval: dto.TimeInterval{
							Start: first,
						},
					},
					{ID: "te-3",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(2)),
						},
					},
				}, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p3",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: api.AllPages(),
				}).Return([]dto.TimeEntry{
					{ID: "te-2",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(1)),
						},
					},
				}, nil)

				return f
			},
			expected: heredoc.Doc(`
				te-1
				te-2
				te-3
			`),
		},
		{
			name: "change timezone",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				tz, _ := time.LoadLocation("America/Sao_Paulo")
				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					TimeZoneLoc:    tz,
					AllowNameForID: false,
				})

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
				rf.Projects = []string{"p"}
				rf.Description = "desc"
				rf.Quiet = true
				rf.Format = `{{ .TimeInterval.Start.Format "` +
					timehlp.FullTimeFormat +
					`" }}`
				return rf
			},
			expected: heredoc.Doc(`
				2006-01-01 22:00:00
				2006-01-04 22:00:00
			`),
		},
		{
			name: "limit number of time entries",
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Limit = 2
				rf.Quiet = true
				return rf
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				f.EXPECT().Config().Return(
					&mocks.SimpleConfig{AllowNameForID: true})

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace: "w",
					UserID:    "u",
					FirstDate: first,
					LastDate:  last,
					PaginationParam: api.PaginationParam{
						Page:     1,
						PageSize: 2,
					},
				}).Return([]dto.TimeEntry{
					{ID: "te-1",
						TimeInterval: dto.TimeInterval{
							Start: first,
						},
					},
					{ID: "te-3",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(2)),
						},
					},
				}, nil)

				return f
			},
			expected: heredoc.Doc(`
				te-1
				te-3
			`),
		},
		{
			name: "limit number of time entries with client filter",
			flags: func(t *testing.T) util.ReportFlags {
				rf := util.NewReportFlags()
				rf.Limit = 2
				rf.Client = "me"
				rf.Quiet = true
				return rf
			},
			factory: func(t *testing.T) cmdutil.Factory {

				f := mocks.NewMockFactory(t)
				f.On("GetUserID").Return("u", nil)
				f.On("GetWorkspaceID").Return("w", nil)

				f.EXPECT().Config().Return(
					&mocks.SimpleConfig{AllowNameForID: true})

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.EXPECT().GetClients(api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{
						{ID: "c1", Name: "me"},
						{ID: "c2", Name: "you"},
					}, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Clients:         []string{"c1"},
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p1", Name: "p1", ClientID: "c1", ClientName: "me"},
					{ID: "p3", Name: "p3", ClientID: "c1", ClientName: "me"},
				}, nil)

				p := api.PaginationParam{Page: 1, PageSize: 2}
				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p1",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: p,
				}).Return([]dto.TimeEntry{
					{ID: "te-1",
						TimeInterval: dto.TimeInterval{
							Start: first,
						},
					},
					{ID: "te-3",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(2)),
						},
					},
				}, nil)

				c.EXPECT().LogRange(api.LogRangeParam{
					Workspace:       "w",
					UserID:          "u",
					ProjectID:       "p3",
					FirstDate:       first,
					LastDate:        last,
					PaginationParam: p,
				}).Return([]dto.TimeEntry{
					{ID: "te-2",
						TimeInterval: dto.TimeInterval{
							Start: first.Add(time.Duration(1)),
						},
					},
				}, nil)

				return f
			},
			expected: heredoc.Doc(`
				te-2
				te-3
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
