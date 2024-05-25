package timehlp_test

import (
	"fmt"
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	now := timehlp.Today()
	nowStr := now.Format("2006-01-02")

	t.Run("now", func(t *testing.T) {
		// this case is special because it is not deterministic
		parsed, err := timehlp.ConvertToTime("now")

		assert.Nil(t, err)
		assert.Equal(t, nowStr, parsed.Format("2006-01-02"))

	})

	tts := []struct {
		name     string
		expected string
		toParse  string
	}{
		{name: "FullTimeFormat", expected: "09:59:01", toParse: fmt.Sprintf("%s %s", nowStr, "09:59:01")},
		{name: "SimplerTimeFormat", expected: "09:59:00", toParse: fmt.Sprintf("%s %s", nowStr, "09:59")},
		{name: "OnlyTimeFormat", expected: "16:03:02", toParse: "16:03:02"},
		{name: "SimplerOnlyTimeFormat", expected: "16:03:00", toParse: "16:03"},
		{name: "SimplerOnlyTimeFormat", expected: "06:03:00", toParse: "06:03"},
		{name: "SimplerOnlyTimeFormatWL", expected: "06:03:00", toParse: "6:03"},
		{name: "SimplestOnlyTimeFormat", expected: "16:03:00", toParse: "1603"},
		{name: "SimplestOnlyTimeFormatWL", expected: "06:03:00", toParse: "603"},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {

			parsed, err := timehlp.ConvertToTime(tt.toParse)

			assert.Nil(t, err)
			assert.Equal(t, fmt.Sprintf("%s %s", nowStr, tt.expected), parsed.Format("2006-01-02 15:04:05"))
		})
	}
}

func TestFailParseTime(t *testing.T) {
	_, err := timehlp.ConvertToTime("2024-05-25 25:61")
	assert.Error(t, err,
		"parsing time \"2024-05-25 25:61:00\": hour out of range")
}
