package list_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestCmdList(t *testing.T) {
	shouldCall := func(t *testing.T) func(
		io.Writer, *util.OutputFlags, []dto.Project) error {
		called := false
		t.Cleanup(func() { assert.True(t, called) })
		return func(w io.Writer, of *util.OutputFlags, p []dto.Project) error {
			called = true
			return nil
		}
	}
	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) cmdutil.Factory
		report  func(*testing.T) func(
			io.Writer, *util.OutputFlags, []dto.Project) error
		err string
	}{
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "archived or not-archived",
			args: []string{"--archived", "--not-archived"},
			err:  "flags can't be used together.*archived.*not-archived",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			args: []string{"-n=a"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("", errors.New("workspace error"))
				return f
			},
		},
		{
			name: "client error",
			err:  "client error",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(nil, errors.New("client error"))
				return f
			},
		},
		{
			name: "lookup client",
			err:  "no client",
			args: []string{"--clients=rockr"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{}, errors.New("no client"))
				return f
			},
		},
		{
			name: "client not found",
			err:  "No client with id or name containing 'other' was found",
			args: []string{"--clients=rockr", "--clients=other"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{{ID: "c1", Name: "Coderockr"}}, nil)
				return f
			},
		},
		{
			name: "http error",
			err:  "http error",
			args: []string{"-n=error"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)
				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Name:            "error",
					Clients:         []string{},
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{}, errors.New("http error"))
				return f
			},
		},
		{
			name: "archived",
			args: []string{
				"--name=cli",
				"--clients=rockr", "--clients", "other",
				"--archived",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{
						{ID: "c1", Name: "Coderockr"},
						{ID: "c2", Name: "Other"},
					}, nil)

				b := true
				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Name:            "cli",
					Clients:         []string{"c1", "c2"},
					Archived:        &b,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{}, nil)
				return f
			},
			report: shouldCall,
		},
		{
			name: "not archived",
			args: []string{
				"--name=cli",
				"--not-archived",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				b := false
				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Name:            "cli",
					Clients:         []string{},
					Archived:        &b,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{}, nil)
				return f
			},
			report: shouldCall,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := func(io.Writer, *util.OutputFlags, []dto.Project) error {
				assert.Fail(t, "failed")
				return nil
			}

			if tt.report != nil {
				r = tt.report(t)
			}

			cmd := list.NewCmdList(tt.factory(t), r)
			cmd.SilenceUsage = true
			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()
			if tt.err == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Regexp(t, tt.err, err.Error())
		})
	}
}

func TestCmdListReport(t *testing.T) {
	pr := []dto.Project{{Name: "Coderockr"}}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, []dto.Project)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Project) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Project) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Project) {
				assert.Equal(t, "{{.ID}}", of.Format)
			},
		},
		{
			name: "report csv",
			args: []string{"--csv"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Project) {
				assert.True(t, of.CSV)
			},
		},
		{
			name: "report default",
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Project) {
				assert.False(t, of.CSV)
				assert.False(t, of.JSON)
				assert.False(t, of.Quiet)
				assert.True(t, of.Format == "")
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)
			c := mocks.NewMockClient(t)
			f.On("Client").Return(c, nil)
			f.On("GetWorkspaceID").
				Return("w", nil)

			c.On("GetProjects", api.GetProjectsParam{
				Workspace:       "w",
				Clients:         []string{},
				PaginationParam: api.AllPages(),
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := list.NewCmdList(f, func(
				_ io.Writer, of *util.OutputFlags, u []dto.Project) error {
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
