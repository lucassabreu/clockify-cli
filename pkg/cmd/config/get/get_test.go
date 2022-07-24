package get_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/get"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newCmd(f cmdutil.Factory) *cobra.Command {
	cmd := get.NewCmdGet(
		f,
		cmdcompl.ValidArgsMap{},
	)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	return cmd
}

func TestGetCmdArgs(t *testing.T) {
	tcs := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "none",
			args: []string{},
			err:  errors.New("requires arg param"),
		},
		{
			name: "two",
			args: []string{"param1", "param2"},
			err:  errors.New("accepts 1 arg(s), received 2"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newCmd(mocks.NewMockFactory(t))
			cmd.SetArgs(tc.args)
			_, err := cmd.ExecuteC()
			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), tc.err.Error())
			}
		})
	}
}

func TestGetCmdRun(t *testing.T) {
	tcs := []struct {
		name   string
		args   []string
		config func(*testing.T) cmdutil.Config
		output string
		err    error
	}{
		{
			name: "token with default format",
			args: []string{"token"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "token").Once().Return("<token-value>")
				return c
			},
			output: "<token-value>\n",
		},
		{
			name: "token with json format",
			args: []string{"token", "--format=json"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "token").Once().Return("token-value")
				return c
			},
			output: `"token-value"`,
		},
		{
			name: "workdays default format",
			args: []string{"workdays"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "workdays").Once().Return([]string{
					"monday",
					"tuesday",
					"sunday",
				})
				return c
			},
			output: heredoc.Doc(`
			- monday
			- tuesday
			- sunday
			`),
		},
		{
			name: "workdays json format",
			args: []string{"workdays", "--format", "json"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "workdays").Once().Return([]string{
					"monday",
					"tuesday",
					"sunday",
				})
				return c
			},
			output: `["monday","tuesday","sunday"]`,
		},
		{
			name: "user.id default format",
			args: []string{"user.id"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "user.id").Once().Return("someuserid")
				return c
			},
			output: "someuserid\n",
		},
		{
			name: "user default format",
			args: []string{"user"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "user").Once().Return(map[string]string{
					"id": "someuserid",
				})
				return c
			},
			output: "id: someuserid\n",
		},
		{
			name: "user json format",
			args: []string{"user", "-f=JSON"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "user").Once().Return(map[string]string{
					"id": "someuserid",
				})
				return c
			},
			output: `{"id":"someuserid"}`,
		},
		{
			name: "invalid format",
			args: []string{"user", "--format", "tmol"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("Get", "user").Return(map[string]string{
					"id": "someuserid",
				})
				return c
			},
			output: ``,
			err:    errors.New("invalid format"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)
			f.On("Config").Once().Return(tc.config(t))

			cmd := newCmd(f)
			cmd.SetArgs(tc.args)
			out := cmd.OutOrStdout().(*bytes.Buffer)
			_, err := cmd.ExecuteC()

			assert.Equal(t, tc.output, out.String())
			if err == nil && assert.NoError(t, err) {
				return
			}

			if !assert.Error(t, err) {
				return
			}
			assert.Equal(t, err.Error(), tc.err.Error())
		})
	}
}
