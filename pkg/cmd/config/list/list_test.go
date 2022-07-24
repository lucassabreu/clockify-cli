package list_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/list"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	tts := []struct {
		name           string
		args           []string
		config         func(t *testing.T) cmdutil.Config
		expectedOutput string
		err            error
	}{
		{
			name:   "no args",
			args:   []string{"param"},
			err:    errors.New(`unknown command "param" for "list"`),
			config: func(t *testing.T) cmdutil.Config { return nil },
		},
		{
			name: "default format",
			args: []string{},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("All").Once().Return(map[string]interface{}{
					"token": "value",
					"user":  map[string]string{"id": "user.id"},
				})
				return c
			},
			expectedOutput: heredoc.Doc(`
			token: value
			user:
			    id: user.id
			`),
		},
		{
			name: "json format",
			args: []string{"--format=json"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("All").Once().Return(map[string]interface{}{
					"token": "value",
					"user":  map[string]string{"id": "user.id"},
				})
				return c
			},
			expectedOutput: `{"token":"value","user":{"id":"user.id"}}`,
		},
		{
			name: "invalid format",
			args: []string{"--format=tmol"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("All").Once().Return(map[string]interface{}{})
				return c
			},
			err: errors.New("invalid format"),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)
			if c := tt.config(t); c != nil {
				f.On("Config").Return(c)
			}

			cmd := list.NewCmdList(f)
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)

			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()
			if tt.err != nil && assert.Error(t, err) {
				assert.EqualError(t, err, tt.err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, b.String())
		})
	}
}
