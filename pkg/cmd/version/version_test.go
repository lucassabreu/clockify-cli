package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/version"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	type c struct {
		name     string
		version  cmdutil.Version
		expected string
	}
	cases := []c{
		c{
			name:     "when no version",
			version:  cmdutil.Version{},
			expected: "Version: , Commit: , Build At:",
		},
		c{
			name: "when default version",
			version: cmdutil.Version{
				Tag: "dev", Commit: "none", Date: "unknown"},
			expected: "Version: dev, Commit: none, Build At: unknown",
		},
		c{
			name: "with valid version",
			version: cmdutil.Version{
				Tag:    "1.0.0",
				Commit: "63595df3d0dd6e4eef2e973631013ca9e06928ef",
				Date:   "2022-07-01T04:10:47Z"},
			expected: "Version: 1.0.0, " +
				"Commit: 63595df3d0dd6e4eef2e973631013ca9e06928ef, " +
				"Build At: 2022-07-01T04:10:47Z",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)
			f.On("Version").Return(tt.version)

			cmd := version.NewCmdVersion(f)
			cmd.SetArgs([]string{})

			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)

			_, err := cmd.ExecuteC()

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, strings.TrimSpace(b.String()))
		})
	}
}
