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

func TestTimeEntriesMarkdownPrint(t *testing.T) {
	t65Min1SecAgo, _ := timehlp.ConvertToTime("-65m1s")
	start, _ := time.Parse(timehlp.FullTimeFormat, "2024-06-15 10:00:01")
	end := start.Add(2*time.Minute + 1*time.Second)

	tts := []struct {
		name   string
		tes    []dto.TimeEntry
		output string
	}{
		{
			name: "open without tags or project",
			tes: []dto.TimeEntry{{
				WorkspaceID:  "w1",
				ID:           "te1",
				Billable:     false,
				Description:  "Open and without project",
				TimeInterval: dto.NewTimeInterval(t65Min1SecAgo, nil),
			}},
			output: heredoc.Docf(`
				## _Time Entry_: te1

				_Time and date_  
				**1:05:01** | Start Time: _%s_ 🗓 Today

				|               |                          |
				|---------------|--------------------------|
				| _Description_ | Open and without project |
				| _Project_     | No Project               |
				| _Tags_        | No Tags                  |
				| _Billable_    | No                       |
			`, t65Min1SecAgo.UTC().Format(timehlp.SimplerOnlyTimeFormat)),
		},
		{
			name: "closed without tags or project",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    false,
				Description: "Closed and without project",
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                            |
				|---------------|----------------------------|
				| _Description_ | Closed and without project |
				| _Project_     | No Project                 |
				| _Tags_        | No Tags                    |
				| _Billable_    | No                         |
			`),
		},
		{
			name: "Closed with project",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    false,
				Description: "With project",
				Project: &dto.Project{
					Name: "Project Name",
				},
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                  |
				|---------------|------------------|
				| _Description_ | With project     |
				| _Project_     | **Project Name** |
				| _Tags_        | No Tags          |
				| _Billable_    | No               |
			`),
		},
		{
			name: "Closed with project with client",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    true,
				Description: "With project",
				Project: &dto.Project{
					Name:       "Project Name",
					ClientName: "Client Name",
				},
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                                |
				|---------------|--------------------------------|
				| _Description_ | With project                   |
				| _Project_     | **Project Name** - Client Name |
				| _Tags_        | No Tags                        |
				| _Billable_    | Yes                            |
			`),
		},
		{
			name: "Closed with project, client and task",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    true,
				Description: "With project",
				Project: &dto.Project{
					Name:       "Project Name",
					ClientName: "Client Name",
				},
				Task: &dto.Task{
					Name: "Task Name",
				},
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                             |
				|---------------|-----------------------------|
				| _Description_ | With project                |
				| _Project_     | **Project Name**: Task Name |
				| _Tags_        | No Tags                     |
				| _Billable_    | Yes                         |
			`),
		},
		{
			name: "Closed with project, client, task and a tag",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    true,
				Description: "With project",
				Project: &dto.Project{
					Name:       "Project Name",
					ClientName: "Client Name",
				},
				Task: &dto.Task{
					Name: "Task Name",
				},
				Tags: []dto.Tag{
					{Name: "Stand-up Meeting"},
				},
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                             |
				|---------------|-----------------------------|
				| _Description_ | With project                |
				| _Project_     | **Project Name**: Task Name |
				| _Tags_        | Stand-up Meeting            |
				| _Billable_    | Yes                         |
			`),
		},
		{
			name: "Closed with project, client, task and tags",
			tes: []dto.TimeEntry{{
				WorkspaceID: "w1",
				ID:          "te1",
				Billable:    true,
				Description: "With project",
				Project: &dto.Project{
					Name:       "Project Name",
					ClientName: "Client Name",
				},
				Task: &dto.Task{
					Name: "Task Name",
				},
				Tags: []dto.Tag{
					{Name: "A Tag with long name"},
					{Name: "Normal tag"},
				},
				TimeInterval: dto.NewTimeInterval(
					start,
					&end,
				),
			}},
			output: heredoc.Doc(`
				## _Time Entry_: te1

				_Time and date_  
				**0:02:01** | 10:00 - 10:02 🗓 06/15/2024

				|               |                                  |
				|---------------|----------------------------------|
				| _Description_ | With project                     |
				| _Project_     | **Project Name**: Task Name      |
				| _Tags_        | A Tag with long name, Normal tag |
				| _Billable_    | Yes                              |
			`),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &strings.Builder{}
			err := timeentry.TimeEntriesMarkdownPrint(tt.tes, buffer)

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.output+"\n", buffer.String())
		})
	}
}
