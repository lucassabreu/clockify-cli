package get_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/get"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

var bFalse = false

func TestCmdGet(t *testing.T) {
	shouldCall := func(t *testing.T) func(
		io.Writer, *util.OutputFlags, dto.Project) error {
		called := false
		t.Cleanup(func() { assert.True(t, called) })
		return func(w io.Writer, of *util.OutputFlags, p dto.Project) error {
			called = true
			return nil
		}
	}
	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) cmdutil.Factory
		report  func(*testing.T) func(
			io.Writer, *util.OutputFlags, dto.Project) error
		err string
	}{
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j", "p1"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			args: []string{"p1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().
					Return("", errors.New("workspace error"))
				return f
			},
		},
		{
			name: "client error",
			err:  "client error",
			args: []string{"p1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)
				f.EXPECT().Client().Return(nil, errors.New("client error"))
				return f
			},
		},
		{
			name: "http error",
			err:  "http error",
			args: []string{"p1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(false)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p1",
				}).
					Return(nil, errors.New("http error"))
				return f
			},
		},
		{
			name: "by id",
			args: []string{"p1"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(false)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p1",
				}).
					Return(&dto.Project{}, nil)
				return f
			},
			report: shouldCall,
		},
		{
			name: "by name",
			args: []string{
				"project",
			},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(true)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{{Name: "project", ID: "p1"}}, nil)

				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p1",
				}).
					Return(&dto.Project{Name: "project", ID: "p1"}, nil)

				return f
			},
			report: shouldCall,
		},
		{
			name: "hydrated",
			args: []string{"-H", "project"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().
					Return("w", nil)

				cf := mocks.NewMockConfig(t)
				f.EXPECT().Config().Return(cf)
				cf.EXPECT().IsAllowNameForID().Return(true)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{{Name: "project", ID: "p1"}}, nil)

				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p1",
					Hydrate:   true,
				}).
					Return(&dto.Project{
						Name: "project", ID: "p1", Hydrated: true},
						nil)
				return f
			},
			report: shouldCall,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := func(io.Writer, *util.OutputFlags, dto.Project) error {
				assert.Fail(t, "failed")
				return nil
			}

			if tt.report != nil {
				r = tt.report(t)
			}

			cmd := get.NewCmdGet(tt.factory(t), r)
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

func TestCmdGetReport(t *testing.T) {
	pr := dto.Project{Name: "Coderockr"}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, dto.Project)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Project) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Project) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Project) {
				assert.Equal(t, "{{.ID}}", of.Format)
			},
		},
		{
			name: "report csv",
			args: []string{"--csv"},
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Project) {
				assert.True(t, of.CSV)
			},
		},
		{
			name: "report default",
			assert: func(t *testing.T, of *util.OutputFlags, _ dto.Project) {
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
			f.EXPECT().GetWorkspaceID().Return("w", nil)

			cf := mocks.NewMockConfig(t)
			f.EXPECT().Config().Return(cf)
			cf.EXPECT().IsAllowNameForID().Return(false)

			c.EXPECT().GetProject(api.GetProjectParam{
				Workspace: "w",
				ProjectID: "p1",
			}).
				Return(&pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := get.NewCmdGet(f, func(
				_ io.Writer, of *util.OutputFlags, u dto.Project) error {
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "p1"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
