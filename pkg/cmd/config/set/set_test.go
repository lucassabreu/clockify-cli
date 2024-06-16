package set_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestSetCmdArgs(t *testing.T) {
	tt := map[string][]string{
		"zero":  {},
		"one":   {"param"},
		"three": {"param", "value", "other value"},
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
	ts := []struct {
		name   string
		args   []string
		config func(t *testing.T) cmdutil.Config
	}{
		{
			name: "set token",
			args: []string{"token", "some value"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("SetString", "token", "some value").Return(nil).Once()
				c.On("Save").Once().Return(nil)
				return c
			},
		},
		{
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
		{
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
		{
			name: "set show client",
			args: []string{cmdutil.CONF_SHOW_CLIENT, "true"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.On("SetString", cmdutil.CONF_SHOW_CLIENT,
					"true").
					Return(nil).Once()
				c.On("Save").Once().Return(nil)
				return c
			},
		},
		{
			name: "set language",
			args: []string{cmdutil.CONF_LANGUAGE, "pt-br"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.EXPECT().SetLanguage(language.BrazilianPortuguese).Once()
				c.EXPECT().Save().Once().Return(nil)
				return c
			},
		},
		{
			name: "set language (iso 639)",
			args: []string{cmdutil.CONF_LANGUAGE, "pt"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				c.EXPECT().SetLanguage(language.Portuguese).Once()
				c.EXPECT().Save().Once().Return(nil)
				return c
			},
		},
		{
			name: "set timezone",
			args: []string{cmdutil.CONF_TIMEZONE, "America/Sao_Paulo"},
			config: func(t *testing.T) cmdutil.Config {
				c := mocks.NewMockConfig(t)
				tz, _ := time.LoadLocation("America/Sao_Paulo")
				c.EXPECT().SetTimeZone(tz).Once()
				c.EXPECT().Save().Once().Return(nil)
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

func TestSetCmdShouldFail(t *testing.T) {
	ts := []struct {
		name string
		args []string
		err  string
	}{
		{
			name: "set language",
			args: []string{cmdutil.CONF_LANGUAGE, "klingon"},
			err:  "klingon is not a valid language.*",
		},
		{
			name: "set timezone",
			args: []string{cmdutil.CONF_TIMEZONE, "Murica"},
			err:  "Murica is not a valid timezone.*",
		},
		{
			name: "set timezone no caps",
			args: []string{cmdutil.CONF_TIMEZONE, "america/sao_paulo"},
			err:  `america/sao_paulo is not a valid timezone`,
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {

			f := mocks.NewMockFactory(t)
			f.EXPECT().Config().Return(mocks.NewMockConfig(t))
			cmd := set.NewCmdSet(f, cmdcompl.ValidArgsMap{})

			b := bytes.NewBufferString("")
			cmd.SetArgs(tc.args)
			cmd.SetErr(b)
			cmd.SetOut(b)
			_, err := cmd.ExecuteC()

			if !assert.Error(t, err) {
				return
			}

			assert.Regexp(t, tc.err, err.Error())
		})
	}
}
