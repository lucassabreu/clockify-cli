package edit_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/edit"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, []dto.Project) error

func TestEditCmd(t *testing.T) {
	tts := []struct {
		name   string
		args   []string
		err    string
		params func(*testing.T) (cmdutil.Factory, report)
	}{
		{
			name: "project is required",
			err:  "requires arg project",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "can only change a project name",
			err:  "`--name` can't be changed for multiple projects",
			args: []string{"cli", "edit", "-n=wrong"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j", "cli"},
			err:  "flags can't be used together.*format.*json.*quiet",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "billable or not",
			args: []string{"--billable", "--not-billable", "cli"},
			err:  "flags can't be used together.*billable.*not-billable",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "active or archived",
			args: []string{"--active", "--archived", "cli"},
			err:  "flags can't be used together.*active.*archived",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "client and no client",
			args: []string{"--client=myself", "--no-client", "cli"},
			err:  "flags can't be used together.*client.*no-client",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "public or private",
			args: []string{"--private", "--public", "cli"},
			err:  "flags can't be used together.*private.*public",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), nil
			},
		},
		{
			name: "workspace error",
			err:  "error",
			args: []string{"cli"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("", errors.New("error"))
				return f, nil
			},
		},
		{
			name: "client error",
			err:  "error",
			args: []string{"cli"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)
				f.On("Client").Return(nil, errors.New("error"))
				return f, nil
			},
		},
		{
			name: "lookup project error",
			err:  "No project with id or name",
			args: []string{"cli", "second"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				cf.On("IsAllowNameForID").Return(true)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{}, nil)

				return f, nil
			},
		},
		{
			name: "fail to update second",
			args: []string{"cli", "second",
				"--public",
				"--billable",
				"--archived",
				"--no-client"},
			err: "error",
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				cf.On("IsAllowNameForID").Return(false)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				b := true
				client := ""
				c.On("UpdateProject", api.UpdateProjectParam{
					Workspace: "w",
					ID:        "cli",
					ClientId:  &client,
					Public:    &b,
					Billable:  &b,
					Archived:  &b,
				}).Return(dto.Project{}, nil)

				c.On("UpdateProject", api.UpdateProjectParam{
					Workspace: "w",
					ID:        "second",
					ClientId:  &client,
					Public:    &b,
					Billable:  &b,
					Archived:  &b,
				}).Return(dto.Project{}, errors.New("error"))

				return f, nil
			},
		},
		{
			name: "update projects",
			args: []string{"cli", "second",
				"--private",
				"--not-billable",
				"--active",
				"--note=active, but not billable",
				"--client=myself"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				cf.On("IsAllowNameForID").Return(true)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p-1", Name: "Clockify CLI"},
					{ID: "p-2", Name: "Second"},
				}, nil)

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Client{
					{ID: "c-1", Name: "Myself"},
				}, nil)

				b := false
				client := "c-1"
				n := "active, but not billable"
				c.On("UpdateProject", api.UpdateProjectParam{
					Workspace: "w",
					ID:        "p-1",
					ClientId:  &client,
					Public:    &b,
					Billable:  &b,
					Archived:  &b,
					Note:      &n,
				}).Return(dto.Project{ID: "cli"}, nil)

				c.On("UpdateProject", api.UpdateProjectParam{
					Workspace: "w",
					ID:        "p-2",
					ClientId:  &client,
					Public:    &b,
					Billable:  &b,
					Archived:  &b,
					Note:      &n,
				}).Return(dto.Project{ID: "edit"}, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					w io.Writer, of *util.OutputFlags, p []dto.Project) error {
					called = true
					assert.Len(t, p, 2)
					return nil
				}
			},
		},
		{
			name: "change name and color",
			args: []string{"first",
				"--name=First Project",
				"--client=myself",
				"--color=0f0"},
			params: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").Return("w", nil)

				cf := mocks.NewMockConfig(t)
				cf.On("IsAllowNameForID").Return(true)
				f.On("Config").Return(cf)

				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetProjects", api.GetProjectsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Project{
					{ID: "p-1", Name: "First"},
				}, nil)

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).Return([]dto.Client{
					{ID: "c-1", Name: "Myself"},
				}, nil)

				client := "c-1"
				c.On("UpdateProject", api.UpdateProjectParam{
					Workspace: "w",
					ID:        "p-1",
					ClientId:  &client,
					Name:      "First Project",
					Color:     "0f0",
				}).Return(dto.Project{ID: "first"}, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					w io.Writer, of *util.OutputFlags, p []dto.Project) error {
					called = true
					assert.Len(t, p, 1)
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			f, r := tt.params(t)
			if r == nil {
				r = func(io.Writer, *util.OutputFlags, []dto.Project) error {
					t.Error("should not be called")
					return nil
				}
			}

			cmd := edit.NewCmdEdit(f, r)
			cmd.SilenceUsage = true
			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()
			if tt.err == "" {
				assert.NoError(t, err)
				return
			}

			if !assert.Error(t, err) {
				return
			}
			assert.Regexp(t, tt.err, err.Error())
		})

	}
}

func TestEditCmdReport(t *testing.T) {
	pr := dto.Project{Name: "Coderockr"}
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

			cf := mocks.NewMockConfig(t)
			cf.On("IsAllowNameForID").Return(false)
			f.On("Config").Return(cf)

			c.On("UpdateProject", api.UpdateProjectParam{
				Workspace: "w",
				ID:        "p-1",
				Name:      "Myself",
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := edit.NewCmdEdit(f, func(
				_ io.Writer, of *util.OutputFlags, u []dto.Project) error {
				called = true
				assert.Contains(t, u, pr)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "-n=Myself", "p-1"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
