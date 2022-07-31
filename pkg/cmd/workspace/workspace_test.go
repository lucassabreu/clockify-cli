package workspace_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/workspace"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestCmdWorkspaces(t *testing.T) {
	tts := []struct {
		name     string
		args     []string
		factory  func(*testing.T) cmdutil.Factory
		err      string
		expected string
	}{
		{
			name: "only_format_or_quiet",
			args: []string{"-q", "--format", "{}"},
			factory: func(t *testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
			err: "the following flags can't be used together: " +
				"`format` and `quiet`",
		},
		{
			name: "invalid_client",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.On("Client").Return(nil, errors.New("no client"))
				return f
			},
			err: "no client",
		},
		{
			name: "quiet",
			args: []string{"-q"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetWorkspaces", api.GetWorkspaces{}).Return(
					[]dto.Workspace{
						{ID: "w1"},
						{ID: "w2"},
					},
					nil,
				)

				return f
			},
			expected: heredoc.Doc(`
				w1
				w2
			`),
		},
		{
			name: "failed to query",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetWorkspaces", api.GetWorkspaces{}).
					Return(nil, errors.New("failed querying"))

				return f
			},
			err: "failed querying",
		},
		{
			name: "format",
			args: []string{"--format", "ID: {{.ID}} | Name: {{ .Name }}"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetWorkspaces", api.GetWorkspaces{}).Return(
					[]dto.Workspace{
						{ID: "w1", Name: "first"},
						{ID: "w2", Name: "last"},
					},
					nil,
				)

				return f
			},
			expected: heredoc.Doc(`
				ID: w1 | Name: first
				ID: w2 | Name: last
			`),
		},
		{
			name: "default",
			args: []string{"--name", "second"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				c.On("GetWorkspaces", api.GetWorkspaces{
					Name: "second",
				}).Return(
					[]dto.Workspace{
						{ID: "w1", Name: "first"},
						{ID: "w2", Name: "last"},
					},
					nil,
				)

				f.On("GetWorkspaceID").Return("first", nil)

				return f
			},
			expected: ".*first.*\n.*last.*",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			cmd := workspace.NewCmdWorkspace(tt.factory(t))
			b := bytes.NewBufferString("")

			cmd.SetOut(b)

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)

				return
			}

			assert.NoError(t, err)
			assert.Regexp(t, tt.expected, b.String())
		})
	}
}
