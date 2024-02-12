package del_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	del "github.com/lucassabreu/clockify-cli/pkg/cmd/task/delete"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, dto.Task) error

var bFalse = false

func TestCmdDelete(t *testing.T) {
	tts := []struct {
		name   string
		err    string
		args   []string
		params func(*testing.T) (
			cmdutil.Factory,
			report,
		)
	}{
		{
			name: "task is required",
			args: []string{},
			err:  "requires arg task",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "project is required",
			err:  "flag.*project.*not set",
			args: []string{"task-id"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "workspace error",
			err:  "w error",
			args: []string{"task-id", "-p", "p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("", errors.New("w error"))
				return f, nil
			},
		},
		{
			name: "client error",
			err:  "c error",
			args: []string{"task-id", "-p", "p-1"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)
				f.On("Client").Return(nil, errors.New("c error"))
				return f, nil
			},
		},
		{
			name: "project lookup error",
			err:  "No project with id or name containing.*p-1",
			args: []string{"task-id", "-p", "p-1"},
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
			name: "task lookup error",
			err:  "No task with id or name containing.*task-id",
			args: []string{"task-id", "-p", "p-1"},
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
				}).Return([]dto.Project{{ID: "p-1"}}, nil)

				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Task{}, nil)

				return f, nil
			},
		},
		{
			name: "task delete error",
			err:  "http error",
			args: []string{"delete", "-p", "p-1"},
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
				}).Return([]dto.Project{{ID: "p-1"}}, nil)

				te := dto.Task{
					ID:        "task-id",
					Name:      "Delete Task",
					ProjectID: "p-1",
				}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Task{te}, nil)

				c.On("DeleteTask", api.DeleteTaskParam{
					Workspace: "w",
					ProjectID: "p-1",
					TaskID:    "task-id",
				}).Return(te, errors.New("http error"))

				return f, nil
			},
		},
		{
			name: "task delete ",
			args: []string{"delete", "-p", "p-1"},
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
				}).Return([]dto.Project{{ID: "p-1"}}, nil)

				te := dto.Task{
					ID:        "task-id",
					Name:      "Delete Task",
					ProjectID: "p-1",
				}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Task{te}, nil)

				c.On("DeleteTask", api.DeleteTaskParam{
					Workspace: "w",
					ProjectID: "p-1",
					TaskID:    "task-id",
				}).Return(te, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called) })
				return f, func(
					_ io.Writer, _ *util.OutputFlags, tr dto.Task) error {
					called = true
					assert.Equal(t, te, tr)
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			f, r := tt.params(t)
			if r == nil {
				r = func(w io.Writer, of *util.OutputFlags, _ dto.Task) error {
					t.Error("should not be called")
					return nil
				}
			}
			cmd := del.NewCmdDelete(f, r)
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

func TestCmdDeleteReport(t *testing.T) {
	pr := dto.Task{Name: "Coderockr"}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, dto.Task)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Task) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Task) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Task) {
				assert.Equal(t, "{{.ID}}", of.Format)
			},
		},
		{
			name: "report csv",
			args: []string{"--csv"},
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Task) {
				assert.True(t, of.CSV)
			},
		},
		{
			name: "report default",
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Task) {
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

			c.On("DeleteTask", api.DeleteTaskParam{
				Workspace: "w",
				ProjectID: "p-1",
				TaskID:    "t-1",
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := del.NewCmdDelete(f, func(
				_ io.Writer, of *util.OutputFlags, u dto.Task) error {
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "t-1", "-p=p-1"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
