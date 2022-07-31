package add_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/add"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/client/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestCmdAdd(t *testing.T) {
	tts := []struct {
		name    string
		args    []string
		factory func(*testing.T) cmdutil.Factory
		err     string
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
			name: "http error",
			err:  "http error",
			args: []string{"-n=error"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("GetWorkspaceID").
					Return("w", nil)
				f.On("Client").Return(c, nil)
				c.On("AddClient", api.AddClientParam{
					Workspace: "w",
					Name:      "error",
				}).
					Return(dto.Client{}, errors.New("http error"))
				return f
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			cmd := add.NewCmdAdd(tt.factory(t),
				func(io.Writer, *util.OutputFlags, dto.Client) error {
					t.Error("should not get here")
					return nil
				})
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
	cl := dto.Client{Name: "Coderockr"}
	tts := []struct {
		name   string
		args   []string
		assert func(*testing.T, *util.OutputFlags, dto.Client)
	}{
		{
			name: "report quiet",
			args: []string{"-q"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Client) {
				assert.True(t, of.Quiet)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Client) {
				assert.True(t, of.JSON)
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.ID}}"},
			assert: func(t *testing.T, of *util.OutputFlags, c dto.Client) {
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

			c.On("AddClient", api.AddClientParam{
				Workspace: "w",
				Name:      "rockr",
			}).
				Return(cl, nil)

			called := false
			t.Cleanup(func() { assert.True(t, called, "was not called") })
			cmd := add.NewCmdAdd(f, func(
				_ io.Writer, of *util.OutputFlags, u dto.Client) error {
				called = true
				assert.Equal(t, cl, u)
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
