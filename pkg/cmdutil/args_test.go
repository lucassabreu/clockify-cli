package cmdutil_test

import (
	"errors"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRequiredNamedArgs(t *testing.T) {
	tt := []struct {
		name    string
		posArgs cobra.PositionalArgs
		args    []string
		err     error
	}{
		{
			name:    "req one and none sent",
			posArgs: cmdutil.RequiredNamedArgs("param1"),
			args:    []string{},
			err:     errors.New("requires arg param1"),
		},
		{
			name:    "req two and none sent",
			posArgs: cmdutil.RequiredNamedArgs("param1", "param2"),
			args:    []string{},
			err: errors.New(
				"requires args param1 and param2; 0 of those received"),
		},
		{
			name:    "req three and one sent",
			posArgs: cmdutil.RequiredNamedArgs("param1", "param2", "param3"),
			args:    []string{"param1"},
			err: errors.New(
				"requires args param1, param2 and param3; 1 of those received",
			),
		},
		{
			name:    "req one and one sent",
			posArgs: cmdutil.RequiredNamedArgs("param1"),
			args:    []string{"param1"},
			err:     nil,
		},
		{
			name:    "req two and two sent",
			posArgs: cmdutil.RequiredNamedArgs("param1", "param2"),
			args:    []string{"param1", "param2"},
			err:     nil,
		},
		{
			name:    "req one and two sent",
			posArgs: cmdutil.RequiredNamedArgs("param1"),
			args:    []string{"param1", "param2"},
			err:     nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmd := cobra.Command{
				Args: tc.posArgs,
				RunE: func(cmd *cobra.Command, args []string) error {
					if tc.err != nil {
						t.Fatal("should not get here")
					}
					return nil
				},
			}

			cmd.SetArgs(tc.args)

			_, err := cmd.ExecuteC()

			if tc.err == nil {
				assert.NoError(t, err)
				return
			}

			var flagErr cmdutil.FlagError
			if !assert.Error(t, err) &&
				assert.ErrorAs(t, err, &flagErr) {
				assert.Equal(t, flagErr.Error(), tc.err.Error())
			}
		})
	}
}
