package done_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/done"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

var bFalse = false

func TestCmdDone(t *testing.T) {
	te := dto.Task{ID: "task-id", Name: "Task with ID", ProjectID: "p-1"}
	shouldCall := func(t *testing.T) func(
		io.Writer, *util.OutputFlags, []dto.Task) error {
		called := false
		t.Cleanup(func() { assert.True(t, called) })
		return func(
			w io.Writer, of *util.OutputFlags, tr []dto.Task) error {
			called = true
			assert.Len(t, tr, 1)
			assert.Equal(t, te, tr[0])
			return nil
		}
	}

	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) cmdutil.Factory
		report  func(*testing.T) func(
			io.Writer, *util.OutputFlags, []dto.Task) error
		err string
	}{
		{
			name: "task id required",
			args: []string{"-p=cli"},
			err:  `requires arg task`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "project required",
			args: []string{"task"},
			err:  `"project" not set`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "project not empty",
			args: []string{"task", "-p=          "},
			err:  `project should not be empty`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "task id not empty",
			args: []string{"   ", "-p=cli"},
			err:  `task id/name should not be empty`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "task id not empty (nice try)",
			args: []string{"not-empty", "  ", "-p=cli"},
			err:  `task id/name should not be empty`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j", "-p=OK", "done"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			args: []string{"done", "-p=b"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)

				f.On("Client").Return(nil, errors.New("client error"))
				return f
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			args: []string{"done", "-p=b"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("", errors.New("workspace error"))
				return f
			},
		},
		{
			name: "lookup project",
			err:  "no project",
			args: []string{"done", "-p=cli"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{}, errors.New("no project"))
				return f
			},
		},
		{
			name: "cant find task",
			err:  "No active task with id or name.*done",
			args: []string{"done", "-p=p-1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{{ID: "p-1"}}, nil)

				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Task{}, nil)

				return f
			},
		},
		{
			name: "fail to find",
			err:  "something went wrong",
			args: []string{"task 1", "task 2", "-p=cli"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{{ID: "p-1", Name: "Cli"}}, nil)

				ts := []dto.Task{
					{ID: "t-1", ProjectID: "p-1", Name: "Task 1"},
					{ID: "t-2", ProjectID: "p-1", Name: "Task 2"},
				}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return(ts, nil)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: ts[0].ProjectID,
					TaskID:    ts[0].ID,
				}).
					Return(ts[0], nil)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					ProjectID: ts[0].ProjectID,
					TaskID:    ts[0].ID,
					Name:      ts[0].Name,
					Status:    api.TaskStatusDone,
				}).
					Return(ts[0], nil)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: ts[1].ProjectID,
					TaskID:    ts[1].ID,
				}).
					Return(ts[1], errors.New("something went wrong"))

				return f
			},
		},
		{
			name: "fail second update",
			err:  "something went wrong",
			args: []string{"task 1", "task 2", "-p=cli"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(true)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{{ID: "p-1", Name: "Cli"}}, nil)

				ts := []dto.Task{
					{ID: "t-1", ProjectID: "p-1", Name: "Task 1"},
					{ID: "t-2", ProjectID: "p-1", Name: "Task 2"},
				}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return(ts, nil)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: ts[0].ProjectID,
					TaskID:    ts[0].ID,
				}).
					Return(ts[0], nil)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					ProjectID: ts[0].ProjectID,
					TaskID:    ts[0].ID,
					Name:      ts[0].Name,
					Status:    api.TaskStatusDone,
				}).
					Return(ts[0], nil)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: ts[1].ProjectID,
					TaskID:    ts[1].ID,
				}).
					Return(ts[1], nil)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					ProjectID: ts[1].ProjectID,
					TaskID:    ts[1].ID,
					Name:      ts[1].Name,
					Status:    api.TaskStatusDone,
				}).
					Return(ts[0], errors.New("something went wrong"))

				return f
			},
		},
		{
			name: "done",
			args: []string{te.ID, "-p=p-1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(false)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: te.ProjectID,
					TaskID:    te.ID,
				}).
					Return(te, nil)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					ProjectID: te.ProjectID,
					TaskID:    te.ID,
					Name:      te.Name,
					Status:    api.TaskStatusDone,
				}).
					Return(te, nil)

				return f
			},
			report: shouldCall,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := func(io.Writer, *util.OutputFlags, []dto.Task) error {
				assert.Fail(t, "failed")
				return nil
			}

			if tt.report != nil {
				r = tt.report(t)
			}

			cmd := done.NewCmdDone(tt.factory(t), r)
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

func TestCmdDoneReport(t *testing.T) {
	tasks := []dto.Task{
		{ID: "t-1", ProjectID: "p-1", Name: "Done Report"},
		{ID: "t-2", ProjectID: "p-1", Name: "Done Cmd"},
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

			fn := func(t dto.Task) {
				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: t.ProjectID,
					TaskID:    t.ID,
				}).
					Return(t, nil)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					ProjectID: t.ProjectID,
					TaskID:    t.ID,
					Name:      t.Name,
					Status:    api.TaskStatusDone,
				}).
					Return(t, nil)

			}

			fn(tasks[0])
			fn(tasks[1])

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := done.NewCmdDone(f, func(
				_ io.Writer, of *util.OutputFlags, l []dto.Task) error {
				called = true
				assert.Contains(t, l, tasks[0])
				assert.Contains(t, l, tasks[1])
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
