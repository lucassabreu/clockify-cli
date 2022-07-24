package set_test

import (
	"bytes"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSetCmdArgs(t *testing.T) {
	tt := map[string][]string{
		"zero":  []string{},
		"one":   []string{"param"},
		"three": []string{"param", "value", "other value"},
	}

	for name := range tt {
		t.Run(name, func(t *testing.T) {
			cmd := set.NewCmdSet(
				mocks.NewMockFactory(t),
				cmdcompl.ValidArgsMap{},
			)
			b := bytes.NewBufferString("")
			cmd.SetArgs(tt[name])
			cmd.SetErr(b)
			cmd.SetOut(b)
			_, err := cmd.ExecuteC()

			assert.Error(t, err)
		})
	}
}

func TestSetCmdRun(t *testing.T) {
	type tc struct {
		name   string
		args   []string
		config func(t *testing.T) cmdutil.Config
	}

	ts := []tc{
		tc{
			name: "set token",
			args: []string{"token", "some value"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("SetString", "token", "some value").Return(nil).Once()
				c.On("Save").Once().Return(nil)
				return c
			},
		},
		tc{
			name: "set weekdays",
			args: []string{cmdutil.CONF_WORKWEEK_DAYS, "SUNDAY,SATURDAY"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("SetStringSlice", cmdutil.CONF_WORKWEEK_DAYS,
					[]string{"sunday", "saturday"}).
					Return(nil).Once()
				c.On("Save").Once().Return(nil)
				return c
			},
		},
		tc{
			name: "set wrong weekdays",
			args: []string{cmdutil.CONF_WORKWEEK_DAYS, "monday,sunday,june"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("SetStringSlice", cmdutil.CONF_WORKWEEK_DAYS,
					[]string{"monday", "sunday"}).
					Return(nil).Once()
				c.On("Save").Once().Return(nil)
				return c
			},
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.config(t)
			f := mocks.NewMockFactory(t)
			f.On("Config").Return(c)
			cmd := set.NewCmdSet(
				f,
				cmdcompl.ValidArgsMap{},
			)
			b := bytes.NewBufferString("")
			cmd.SetArgs(tc.args)
			cmd.SetErr(b)
			cmd.SetOut(b)
			_, err := cmd.ExecuteC()

			assert.NoError(t, err)
		})
	}

}
