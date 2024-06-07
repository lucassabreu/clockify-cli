package timeentry_test

import (
	"strings"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	timeentry "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestTimeEntriesTotalDurationOnlyAsFloat_ShouldUseUserLanguage(
	t *testing.T) {

	start, _ := timehlp.ConvertToTime("2024-01-01 00:00")
	end := start.Add(time.Hour*1000 + (time.Minute * 31))

	tes := []dto.TimeEntry{
		{TimeInterval: dto.NewTimeInterval(start, &end)},
	}

	tts := []struct {
		name     string
		language language.Tag
		output   string
	}{
		{language: language.English, output: "1,000.517"},
		{language: language.German, output: "1.000,517"},
		{language: language.MustParse("pt-br"), output: "1.000,517"},
		{language: language.Spanish, output: "1.000,517"},
		{language: language.Afrikaans, output: "1\u00a0000,517"},
	}

	for _, tt := range tts {
		t.Run(tt.language.String(), func(t *testing.T) {
			buffer := strings.Builder{}

			err := timeentry.TimeEntriesTotalDurationOnlyAsFloat(
				tes, &buffer, tt.language)

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.output+"\n", buffer.String())
		})
	}
}
