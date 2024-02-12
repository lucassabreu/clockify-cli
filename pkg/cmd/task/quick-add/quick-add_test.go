package quickadd_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	quickadd "github.com/lucassabreu/clockify-cli/pkg/cmd/task/quick-add"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/task/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

var bFalse = false

func TestCmdQuickAdd(t *testing.T) {
	shouldCall := func(t *testing.T) func(
		io.Writer, *util.OutputFlags, []dto.Task) error {
		called := false
		t.Cleanup(func() { assert.True(t, called) })
		return func(
			_ io.Writer, _ *util.OutputFlags, ts []dto.Task) error {
			called = true
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
			name: "only one format",
			args: []string{"--format={}", "-q", "-j", "OK", "-p=OK"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "name required",
			args: []string{"-p=OK"},
			err:  `requires arg name`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "project required",
			args: []string{"OK"},
			err:  `"project" not set`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			args: []string{"a", "-p=b"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Client().Return(nil, errors.New("client error"))
				return f
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			args: []string{"a", "-p=b"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Client().Return(mocks.NewMockClient(t), nil)
				f.EXPECT().GetWorkspaceID().
					Return("", errors.New("workspace error"))
				return f
			},
		},
		{
			name: "lookup project",
			err:  "no project",
			args: []string{"error", "-p=cli"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)
				f.EXPECT().Client().Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(true)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{}, errors.New("no project"))
				return f
			},
		},
		{
			name: "http error",
			err:  "http error",
			args: []string{"error", "-p=ok"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)
				f.EXPECT().Client().Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(false)

				c.EXPECT().AddTask(api.AddTaskParam{
					Workspace: "w",
					ProjectID: "ok",
					Name:      "error",
				}).
					Return(dto.Task{}, errors.New("http error"))
				return f
			},
		},
		{
			name: "add one task",
			args: []string{
				"--project=cli",
				"Add",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)
				f.EXPECT().Client().Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(true)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return(
					[]dto.Project{{ID: "p-1", Name: "Clockify CLI"}}, nil)

				c.EXPECT().AddTask(api.AddTaskParam{
					Workspace: "w",
					Name:      "Add",
					ProjectID: "p-1",
				}).
					Return(dto.Task{ID: "t-id"}, nil)

				return f
			},
			report: shouldCall,
		},
		{
			name: "add multiple tasks",
			args: []string{
				"--project=cli",
				"Task 00",
				"Task 01",
				"Task 02",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)
				f.EXPECT().Client().Return(c, nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(true)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return(
					[]dto.Project{{ID: "p-1", Name: "Clockify CLI"}}, nil)

				c.EXPECT().AddTask(api.AddTaskParam{
					Workspace: "w",
					Name:      "Task 00",
					ProjectID: "p-1",
				}).
					Return(dto.Task{ID: "00"}, nil)

				c.EXPECT().AddTask(api.AddTaskParam{
					Workspace: "w",
					Name:      "Task 01",
					ProjectID: "p-1",
				}).
					Return(dto.Task{ID: "01"}, nil)

				c.EXPECT().AddTask(api.AddTaskParam{
					Workspace: "w",
					Name:      "Task 02",
					ProjectID: "p-1",
				}).
					Return(dto.Task{ID: "02"}, nil)

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

			cmd := quickadd.NewCmdQuickAdd(tt.factory(t), r)
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

func TestCmdQuickAddReport(t *testing.T) {
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
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().GetWorkspaceID().
				Return("w", nil)

			cf := mocks.NewMockConfig(t)
			f.EXPECT().Config().Return(cf)

			cf.EXPECT().IsAllowNameForID().Return(false)

			c.EXPECT().AddTask(api.AddTaskParam{
				Workspace: "w",
				ProjectID: "p-1",
				Name:      "Task Add",
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := quickadd.NewCmdQuickAdd(f, func(
				_ io.Writer, of *util.OutputFlags, ts []dto.Task) error {
				u := ts[0]
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "Task Add", "-p=p-1"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
