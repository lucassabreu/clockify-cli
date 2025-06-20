package timeentry_test

import (
	"strings"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api/dto"
	timeentry "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
)

func TestTimeEntriesDefaultPrint(t *testing.T) {
	start, _ := time.Parse(timehlp.FullTimeFormat, "2024-06-15 10:00:01")
	end := start.Add(2*time.Minute + 1*time.Second)

	tts := []struct {
		name   string
		opts   timeentry.TimeEntryOutputOptions
		tes    []dto.TimeEntry
		output string
	}{
		{
			name: "show clients on its own column",
			opts: timeentry.TimeEntryOutputOptions{
				ShowClients:       true,
				TimeFormat:        timehlp.FullTimeFormat,
				ShowTotalDuration: true,
			},
			tes: []dto.TimeEntry{
				{
					WorkspaceID: "w1",
					ID:          "dasdasdasdaasdasdasdasda",
					Billable:    true,
					Description: "With project",
					Project: &dto.Project{
						Name:       "Project Name",
						ClientName: "Client Name",
					},
					Tags: []dto.Tag{{
						ID:   "tag1",
						Name: "Tag Name",
					}},
					TimeInterval: dto.NewTimeInterval(
						start,
						&end,
					),
				},
				{
					WorkspaceID: "w1",
					ID:          "dfsdfsdfsdffsdfsdfsdfsdf",
					Billable:    true,
					Description: "Without project",
					Tags: []dto.Tag{{
						ID:   "tag1",
						Name: "Tag Name",
					}},
					TimeInterval: dto.NewTimeInterval(
						start,
						&end,
					),
				},
			},
			output: heredoc.Docf(`
				+--------------------------+---------------------+---------------------+---------+--------------+-------------+-----------------+-----------------+
				|            ID            |        START        |         END         |   DUR   |   PROJECT    |   CLIENT    |   DESCRIPTION   |      TAGS       |
				+--------------------------+---------------------+---------------------+---------+--------------+-------------+-----------------+-----------------+
				| dasdasdasdaasdasdasdasda | %s | %s | 0:02:01 | Project Name | Client Name | With project    | Tag Name (tag1) |
				+--------------------------+---------------------+---------------------+---------+--------------+-------------+-----------------+-----------------+
				| dfsdfsdfsdffsdfsdfsdfsdf | %s | %s | 0:02:01 |              |             | Without project | Tag Name (tag1) |
				+--------------------------+---------------------+---------------------+---------+--------------+-------------+-----------------+-----------------+
				| TOTAL                    |                     |                     | 0:04:02 |              |             |                 |                 |
				+--------------------------+---------------------+---------------------+---------+--------------+-------------+-----------------+-----------------+
				`,
				start.In(time.Local).Format(timehlp.FullTimeFormat), end.In(time.Local).Format(timehlp.FullTimeFormat),
				start.In(time.Local).Format(timehlp.FullTimeFormat), end.In(time.Local).Format(timehlp.FullTimeFormat),
			),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &strings.Builder{}
			err := timeentry.TimeEntriesPrint(tt.opts)(tt.tes, buffer)

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.output, buffer.String())
		})
	}
}
