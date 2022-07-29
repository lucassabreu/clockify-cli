package list_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, []dto.Client) error

func TestCmdList(t *testing.T) {
	defReport := func(io.Writer, *util.OutputFlags, []dto.Client) error {
		return errors.New("should not call")
	}

	cs := []dto.Client{{Name: "Coderockr"}}
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
				return mocks.NewMockFactory(t), defReport
			},
		},
		{
			name: "archived or not",
			args: []string{"--archived", "--not-archived"},
			err:  "flags can't be used together.*archived.*not-archived",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				return mocks.NewMockFactory(t), defReport
			},
		},
		{
			name: "client error",
			err:  "client error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(nil, errors.New("client error"))
				return f, defReport
			},
		},
		{
			name: "workspace error",
			err:  "workspace error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("GetWorkspaceID").
					Return("", errors.New("workspace error"))
				return f, defReport
			},
		},
		{
			name: "http error",
			err:  "http error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)
				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Client{}, errors.New("http error"))
				return f, defReport
			},
		},
		{
			name: "only archived",
			args: []string{"--archived"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				b := true
				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					Archived:        &b,
					PaginationParam: api.AllPages(),
				}).
					Return(cs, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, l []dto.Client) error {
					called = true
					assert.Equal(t, cs, l)
					return nil
				}
			},
		},
		{
			name: "not archived",
			args: []string{"--not-archived"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				b := false
				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					Archived:        &b,
					PaginationParam: api.AllPages(),
				}).
					Return(cs, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, l []dto.Client) error {
					called = true
					assert.Equal(t, cs, l)
					return nil
				}
			},
		},
		{
			name: "report quiet",
			args: []string{"--name=rockr", "-q"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				f.On("GetWorkspaceID").
					Return("w", nil)

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					Name:            "rockr",
					PaginationParam: api.AllPages(),
				}).
					Return(cs, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.Client) error {
					called = true
					assert.Equal(t, cs, u)
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

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return(cs, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.Client) error {
					called = true
					assert.Equal(t, cs, u)
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

				c.On("GetClients", api.GetClientsParam{
					Workspace:       "w",
					PaginationParam: api.AllPages(),
				}).
					Return(cs, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u []dto.Client) error {
					called = true
					assert.Equal(t, cs, u)
					assert.Equal(t, "{{.Name}}", of.Format)
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			cmd := list.NewCmdList(tt.factory(t))
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
