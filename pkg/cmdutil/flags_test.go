package cmdutil_test

import (
	"errors"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	name  string
	param map[string]bool
	err   error
}

func testcases() []testcase {
	return []testcase{
		{
			name: "all false",
			param: map[string]bool{
				"pos1": false,
				"pos2": false,
				"pos3": false,
			},
		},
		{
			name:  "empty",
			param: map[string]bool{},
		},
		{
			name: "pos1 and pos2 are true",
			param: map[string]bool{
				"pos1": true,
				"pos2": true,
				"pos3": false,
			},
			err: errors.New(
				"the following flags can't be used together: " +
					"`pos1` and `pos2`"),
		},
		{
			name: "pos1, pos2 and pos3 are true",
			param: map[string]bool{
				"pos1": true,
				"pos2": true,
				"pos4": false,
				"pos3": true,
			},
			err: errors.New(
				"the following flags can't be used together: " +
					"`pos1`, `pos2` and `pos3`"),
		},
		{
			name: "pos1 and pos4 are true",
			param: map[string]bool{
				"pos1": true,
				"pos2": false,
				"pos3": false,
				"pos4": true,
			},
			err: errors.New(
				"the following flags can't be used together: " +
					"`pos1` and `pos4`"),
		},
	}
}

func TestXorFlag(t *testing.T) {
	for _, tc := range testcases() {
		t.Run(tc.name, func(t *testing.T) {
			err := cmdutil.XorFlag(tc.param)
			if tc.err == nil && assert.NoError(t, err) {
				return
			}

			assert.Error(t, err)
			var fErr *cmdutil.FlagError
			assert.ErrorAs(t, err, &fErr)
			assert.EqualError(t, tc.err, err.Error())
		})
	}
}

func TestXorFlagSet(t *testing.T) {
	for _, tc := range testcases() {
		t.Run(tc.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flags := make([]string, len(tc.param))
			args := []string{}
			for fl := range tc.param {
				flags = append(flags, fl)
				fs.Bool(fl, false, "help")
				if tc.param[fl] {
					args = append(args, "--"+fl)
				}
			}
			_ = fs.Parse(args)

			err := cmdutil.XorFlagSet(fs, flags...)
			if tc.err == nil && assert.NoError(t, err) {
				return
			}

			assert.Error(t, err)
			var fErr *cmdutil.FlagError
			assert.ErrorAs(t, err, &fErr)
			assert.EqualError(t, tc.err, err.Error())
		})
	}
}
