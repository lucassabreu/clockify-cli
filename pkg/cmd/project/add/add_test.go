package add_test

import (
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/add"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/project/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCmdAdd(t *testing.T) {
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
			args: []string{"--format={}", "-q", "-j", "-n=OK"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "random-color or color",
			args: []string{"--color=f00", "--random-color", "-n=OK"},
			err:  "flags can't be used together.*color.*random-color",
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "name required",
			err:  `"name" not set`,
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			args: []string{"-n=a"},
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
			args: []string{"-n=a"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("", errors.New("workspace error"))
				return f
			},
		},
		{
			name: "lookup client",
			err:  "no client",
			args: []string{"-n=error", "--client=rockr"},
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
			name: "http error",
			err:  "http error",
			args: []string{"-n=error"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)
				c.On("AddProject", api.AddProjectParam{
					Workspace: "w",
					Name:      "error",
				}).
					Return(dto.Project{}, errors.New("http error"))
				return f
			},
		},
		{
			name: "add project",
			args: []string{
				"--name=Clockify",
				"--client=self",
				"--color=f00",
				"--note", "This one",
				"--public",
				"--billable",
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

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{{ID: "c-1", Name: "Myself"}}, nil)

				c.On("AddProject", api.AddProjectParam{
					Workspace: "w",
					Name:      "Clockify",
					ClientId:  "c-1",
					Note:      "This one",
					Public:    true,
					Billable:  true,
					Color:     "f00",
				}).
					Return(dto.Project{ID: "project-id"}, nil)

				return f
			},
			report: func(t *testing.T) func(
				io.Writer, *util.OutputFlags, dto.Project) error {
				called := false
				t.Cleanup(func() { assert.True(t, called) })
				return func(
					w io.Writer, of *util.OutputFlags, p dto.Project) error {
					called = true
					assert.Equal(t, "project-id", p.ID)
					return nil
				}
			},
		},
		{
			name: "add with random color",
			args: []string{
				"--name=Clockify",
				"--client=c-id",
				"--random-color",
				"--note", "This one",
				"--public",
				"--billable",
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

				c.On("AddProject", mock.AnythingOfType("api.AddProjectParam")).
					Run(func(args mock.Arguments) {
						p := args.Get(0).(api.AddProjectParam)

						assert.Equal(t, p.Workspace, "w")
						assert.Equal(t, p.Name, "Clockify")
						assert.Equal(t, p.ClientId, "c-id")
						assert.Equal(t, p.Note, "This one")
						assert.Equal(t, p.Public, true)
						assert.Equal(t, p.Billable, true)
						assert.Regexp(t,
							regexp.MustCompile("#[0-9a-f]{6}"), p.Color)
					}).
					Return(dto.Project{ID: "project-id"}, nil)

				return f
			},
			report: func(t *testing.T) func(
				io.Writer, *util.OutputFlags, dto.Project) error {
				called := false
				t.Cleanup(func() { assert.True(t, called) })
				return func(
					w io.Writer, of *util.OutputFlags, p dto.Project) error {
					called = true
					assert.Equal(t, "project-id", p.ID)
					return nil
				}
			},
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

			cmd := add.NewCmdAdd(tt.factory(t), r)
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

func TestCmdAddReport(t *testing.T) {
	pr := dto.Project{Name: "Coderockr"}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, dto.Project)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Project) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Project) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Project) {
				assert.Equal(t, "{{.ID}}", of.Format)
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

			c.On("AddProject", api.AddProjectParam{
				Workspace: "w",
				Name:      "rockr",
			}).
				Return(pr, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := add.NewCmdAdd(f, func(
				_ io.Writer, of *util.OutputFlags, u dto.Project) error {
				called = true
				assert.Equal(t, pr, u)
				tt.assert(t, of, u)
				return nil
			})
			cmd.SilenceUsage = true
			cmd.SetArgs(append(tt.args, "-n=rockr"))

			_, err := cmd.ExecuteC()
			assert.NoError(t, err)
		})
	}
}
