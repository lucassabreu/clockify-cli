package list_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, []dto.Task) error

var bFalse = false

func TestCmdList(t *testing.T) {
	tts := []struct {
		name   string
		err    string
		args   []string
		params func(*testing.T) (cmdutil.Factory, report)
	}{
		{
			name: "missing project",
			err:  "required flag.*project",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "workspace error",
			err:  "error",
			args: []string{"-p=p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("", errors.New("error"))
				return f, nil
			},
		},
		{
			name: "client error",
			err:  "error",
			args: []string{"-p=p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)
				f.On("Client").Return(nil, errors.New("error"))
				return f, nil
			},
		},
		{
			name: "project lookup error",
			err:  "No project with id or name containing.*p-1",
			args: []string{"-p=p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{}, nil)

				return f, nil
			},
		},
		{
			name: "list error",
			err:  "error",
			args: []string{"-p=p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				cf.On("IsAllowNameForID").Return(false)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Task{}, errors.New("error"))

				return f, nil
			},
		},
		{
			name: "list active with name",
			args: []string{"-p=p-1", "--active", "-n=list"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)

				cf.On("IsAllowNameForID").Return(true)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p-1", Name: "Cli"},
				}, nil)

				l := []dto.Task{
					{ID: "1", Name: "List"},
					{ID: "2", Name: "List Tasks"},
				}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					Active:          true,
					Name:            "list",
					PaginationParam: api.AllPages(),
				}).Return(l, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					w io.Writer, of *util.OutputFlags, r []dto.Task) error {
					called = true
					assert.Equal(t, l, r)
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			f, r := tt.params(t)
			if r == nil {
				r = func(
					_ io.Writer, of *util.OutputFlags, l []dto.Task) error {
					t.Error("should not be called")
					return nil
				}
			}
			cmd := list.NewCmdList(f, r)
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
	tasks := []dto.Task{
		{ID: "t-1", ProjectID: "p-1", Name: "List Report"},
		{ID: "t-2", ProjectID: "p-1", Name: "List Cmd"},
	}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, []dto.Task)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Task) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Task) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Task) {
				assert.Equal(t, "{{.ID}}", of.Format)
			},
		},
		{
			name: "report csv",
			args: []string{"--csv"},
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Task) {
				assert.True(t, of.CSV)
			},
		},
		{
			name: "report default",
			assert: func(t *testing.T, of *util.OutputFlags, _ []dto.Task) {
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

			cf := mocks.NewMockConfig(t)
			f.On("Config").Return(cf)

			cf.On("IsAllowNameForID").Return(false)

			c.On("GetTasks", api.GetTasksParam{
				Workspace:       "w",
				ProjectID:       "p-1",
				PaginationParam: api.AllPages(),
			}).
				Return(tasks, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := list.NewCmdList(f, func(
				_ io.Writer, of *util.OutputFlags, l []dto.Task) error {
				called = true
				assert.Equal(t, l, tasks)
				tt.assert(t, of, l)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "t-1", "-p=p-1", "t-2"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
