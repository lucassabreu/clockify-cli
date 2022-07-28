package user_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, []dto.User) error

func TestCmdUser(t *testing.T) {
	defReport := func(t *testing.T) report {
		return func(io.Writer, *util.OutputFlags, []dto.User) error {
			t.Error("should not report users")
			return nil
		}
	}
	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) (cmdutil.Factory, report)
		err     string
	}{
		{
			name: "only one format",
			args: []string{"--format={}", "-q", "-j"},
			err:  "flags can't be used together.*format.*json.*quiet",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), defReport(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("Client").Return(nil, errors.New("client error"))
				return f, defReport(t)
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("Client").Return(mocks.NewMockClient(t), nil)
				f.On("GetWorkspaceID").
					Return("", errors.New("workspace error"))
				return f, defReport(t)
			},
		},
		{
			name: "http error",
			err:  "http error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)
				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.User{}, errors.New("http error"))
				return f, defReport(t)
			},
		},
		{
			name: "report quiet",
			args: []string{"--email=john@due.com", "-q"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				list := []dto.User{{Email: "john@due.com"}}
				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					Email:           "john@due.com",
					PaginationParam: api.AllPages(),
				}).
					Return(list, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.User) error {
					called = true
					assert.Equal(t, list, u)
					assert.True(t, of.Quiet)
					return nil
				}
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				list := []dto.User{{Email: "john@due.com"}}
				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return(list, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.User) error {
					called = true
					assert.Equal(t, list, u)
					assert.True(t, of.JSON)
					return nil
				}
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.Name}}"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				list := []dto.User{{Email: "john@due.com"}}
				c.On("WorkspaceUsers", api.WorkspaceUsersParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return(list, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.User) error {
					called = true
					assert.Equal(t, list, u)
					assert.Equal(t, "{{.Name}}", of.Format)
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			cmd := user.NewCmdUser(tt.factory(t))
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
