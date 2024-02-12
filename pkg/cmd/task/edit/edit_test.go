package edit_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/edit"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

var bFalse = false

func TestCmdEdit(t *testing.T) {
	te := dto.Task{ID: "task-id"}
	shouldCall := func(t *testing.T) func(
		io.Writer, *util.OutputFlags, dto.Task) error {
		called := false
		t.Cleanup(func() { assert.True(t, called) })
		return func(
			w io.Writer, of *util.OutputFlags, tr dto.Task) error {
			called = true
			assert.Equal(t, te, tr)
			return nil
		}
	}

	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) cmdutil.Factory
		report  func(*testing.T) func(
			io.Writer, *util.OutputFlags, dto.Task) error
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
			name: "task id not empty",
			args: []string{"   ", "-p=cli"},
			err:  `task id should not be empty`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j", "-p=OK", "edit"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "billable or not",
			args: []string{"--billable", "--not-billable", "edit", "-p=OK"},
			err:  "flags can't be used together.*billable.*not-billable",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "assignee or no assignee",
			args: []string{"--assignee=l", "--no-assignee", "edit", "-p=OK"},
			err:  "flags can't be used together.*assignee.*no-assignee",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			args: []string{"edit", "-p=b"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(false)

				f.On("Client").Return(nil, errors.New("client error"))
				return f
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			args: []string{"edit", "-p=b"},
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
			args: []string{"edit", "-p=cli"},
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
			name: "lookup user",
			err:  "no user",
			args: []string{"edit", "-p=cli", "-A=who"},
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
					Return([]dto.Project{{Name: "Cli"}}, nil)

				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.User{}, errors.New("no user"))

				return f
			},
		},
		{
			name: "lookup task",
			err:  "No task with id or name.*edit",
			args: []string{"edit", "-p=cli"},
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

				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Task{{Name: "other one"}}, nil)

				return f
			},
		},
		{
			name: "http error",
			err:  "http error",
			args: []string{"task-id", "-p=ok", "-n=new"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(false)

				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace: "w",
					TaskID:    "task-id",
					ProjectID: "ok",
					Name:      "new",
				}).
					Return(dto.Task{}, errors.New("http error"))
				return f
			},
		},
		{
			name: "edit billable task",
			args: []string{
				"edit",
				"--name=Edit",
				"--project=cli",
				"--billable",
				"--no-assignee",
				"--estimate", "32",
			},
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
				}).Return(
					[]dto.Project{{ID: "p-1", Name: "Clockify CLI"}}, nil)

				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return(
					[]dto.Task{{ID: "t-1", Name: "Edit Command"}}, nil)

				b := true
				e := time.Hour * 32
				us := []string{}
				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace:   "w",
					TaskID:      "t-1",
					Name:        "Edit",
					ProjectID:   "p-1",
					Billable:    &b,
					Estimate:    &e,
					AssigneeIDs: &us,
				}).
					Return(te, nil)

				return f
			},
			report: shouldCall,
		},
		{
			name: "edit non-billable task",
			args: []string{
				"edit",
				"--project=p-1",
				"--assignee", "lucas",
				"--assignee=john",
				"--not-billable",
			},
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
				}).Return(
					[]dto.Project{{ID: "p-1", Name: "Clockify CLI"}}, nil)

				ts := dto.Task{ID: "t-1", Name: "Edit Command"}
				c.On("GetTasks", api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p-1",
					PaginationParam: api.AllPages(),
				}).Return(
					[]dto.Task{ts}, nil)

				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.User{
						{ID: "u-1", Name: "Lucas Abreu"},
						{ID: "u-2", Name: "John Due"},
					}, nil)

				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: "p-1",
					TaskID:    "t-1",
				}).Return(ts, nil)

				b := false
				as := []string{"u-1", "u-2"}
				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace:   "w",
					TaskID:      "t-1",
					Name:        "Edit Command",
					ProjectID:   "p-1",
					AssigneeIDs: &as,
					Billable:    &b,
				}).
					Return(te, nil)

				return f
			},
			report: shouldCall,
		},
		{
			name: "edit no allow name for id",
			args: []string{
				"t-1",
				"--project=p-1",
				"--assignee", "u-1",
				"--assignee=u-2",
				"--not-billable",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.On("Config").Return(cf)
				cf.On("IsAllowNameForID").Return(false)

				ts := dto.Task{ID: "t-1", Name: "Edit Command"}
				c.On("GetTask", api.GetTaskParam{
					Workspace: "w",
					ProjectID: "p-1",
					TaskID:    "t-1",
				}).Return(ts, nil)

				b := false
				as := []string{"u-1", "u-2"}
				c.On("UpdateTask", api.UpdateTaskParam{
					Workspace:   "w",
					TaskID:      "t-1",
					Name:        "Edit Command",
					ProjectID:   "p-1",
					AssigneeIDs: &as,
					Billable:    &b,
				}).
					Return(te, nil)

				return f
			},
			report: shouldCall,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := func(io.Writer, *util.OutputFlags, dto.Task) error {
				assert.Fail(t, "failed")
				return nil
			}

			if tt.report != nil {
				r = tt.report(t)
			}

			cmd := edit.NewCmdEdit(tt.factory(t), r)
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

func TestCmdEditReport(t *testing.T) {
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

			c.On("UpdateTask", api.UpdateTaskParam{
				Workspace: "w",
				ProjectID: "p-1",
				TaskID:    "t-1",
				Name:      "Task Edit",
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := edit.NewCmdEdit(f, func(
				_ io.Writer, of *util.OutputFlags, u dto.Task) error {
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "t-1", "-n", "Task Edit", "-p=p-1"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
