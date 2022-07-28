package me_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/me"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/user/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
)

type report func(io.Writer, *util.OutputFlags, dto.User) error

func TestCmdMe(t *testing.T) {
	defReport := func(t *testing.T) report {
		return func(io.Writer, *util.OutputFlags, dto.User) error {
			t.Error("should not report users")
			return nil
		}
	}
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
				return mocks.NewMockFactory(t), defReport(t)
			},
		},
		{
			name: "client error",
			err:  "client error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				f.On("Client").Return(nil, errors.New("client error"))
				return f, defReport(t)
			},
		},
		{
			name: "http error",
			err:  "http error",
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)
				c.On("GetMe").
					Return(dto.User{}, errors.New("http error"))
				return f, defReport(t)
			},
		},
		{
			name: "report json",
			args: []string{"--json"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				r := dto.User{Email: "john@due.com"}
				c.On("GetMe").Return(r, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u dto.User) error {
					called = true
					assert.Equal(t, r, u)
					assert.True(t, of.JSON)
					return nil
				}
			},
		},
		{
			name: "report quiet",
			args: []string{"-q"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				r := dto.User{Email: "john@due.com"}
				c.On("GetMe").Return(r, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u dto.User) error {
					called = true
					assert.Equal(t, r, u)
					assert.True(t, of.Quiet)
					return nil
				}
			},
		},
		{
			name: "report format",
			args: []string{"--format={{.Email}}"},
			factory: func(t *testing.T) (cmdutil.Factory, report) {
				f := mocks.NewMockFactory(t)
				c := mocks.NewMockClient(t)
				f.On("Client").Return(c, nil)

				r := dto.User{Email: "john@due.com"}
				c.On("GetMe").Return(r, nil)

				called := false
				t.Cleanup(func() { assert.True(t, called, "was not called") })
				return f, func(
					_ io.Writer, of *util.OutputFlags, u dto.User) error {
					called = true
					assert.Equal(t, r, u)
					assert.Equal(t, of.Format, "{{.Email}}")
					return nil
				}
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			cmd := me.NewCmdMe(tt.factory(t))
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
