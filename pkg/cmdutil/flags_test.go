package cmdutil_test

import (
	"errors"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestXorFlag(t *testing.T) {
	tt := []struct {
		name  string
		param map[string]bool
		err   error
	}{
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

	for _, tc := range tt {
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
