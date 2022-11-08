package api_test

import (
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	. "github.com/lucassabreu/clockify-cli/internal/testhlp"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
)

func TestCreateTimeEntry(t *testing.T) {
	uri := "/v1/workspaces/" + exampleID + "/time-entries"
	end := MustParseTime(timehlp.SimplerTimeFormat, "2022-11-07 11:00")
	bTrue := true
	bFalse := false

	tts := []testCase{
		&simpleTestCase{
			name:  "workspace is required",
			param: api.CreateTimeEntryParam{},
			err:   "workspace is required",
		},
		&simpleTestCase{
			name: "workspace is valid",
			param: api.CreateTimeEntryParam{
				Workspace: "w",
			},
			err: "workspace .* is not valid ID",
		},
		&simpleTestCase{
			name: "with just start time",
			param: api.CreateTimeEntryParam{
				Workspace: exampleID,
				Start: MustParseTime(timehlp.SimplerTimeFormat,
					"2022-11-07 10:00"),
			},

			requestMethod: "post",
			requestUrl:    uri,
			requestBody:   `{"start":"2022-11-07T10:00:00Z"}`,

			responseStatus: 200,
			responseBody:   `{"id": "1"}`,

			result: dto.TimeEntryImpl{ID: "1"},
		},
		&simpleTestCase{
			name: "with all options (billable)",
			param: api.CreateTimeEntryParam{
				Workspace: exampleID,
				Start: MustParseTime(timehlp.SimplerTimeFormat,
					"2022-11-07 10:00"),
				End:         &end,
				Billable:    &bTrue,
				Description: "new entry",
				ProjectID:   "p",
				TaskID:      "t",
				TagIDs:      []string{"tag1", "tag2"},
			},

			requestMethod: "post",
			requestUrl:    uri,
			requestBody: `{
				"start":"2022-11-07T10:00:00Z",
				"end":"2022-11-07T11:00:00Z",
				"billable": true,
				"description": "new entry",
				"projectId": "p",
				"taskId": "t",
				"tagIds": ["tag1","tag2"]
			}`,

			responseStatus: 200,
			responseBody:   `{"id": "1"}`,

			result: dto.TimeEntryImpl{ID: "1"},
		},
		&simpleTestCase{
			name: "not billable",
			param: api.CreateTimeEntryParam{
				Workspace: exampleID,
				Start: MustParseTime(timehlp.SimplerTimeFormat,
					"2022-11-07 10:00"),
				Billable:    &bFalse,
				Description: "new entry",
				ProjectID:   "p",
			},

			requestMethod: "post",
			requestUrl:    uri,
			requestBody: `{
				"start":"2022-11-07T10:00:00Z",
				"billable": false,
				"description": "new entry",
				"projectId": "p"
			}`,

			responseStatus: 200,
			responseBody:   `{"id": "1"}`,

			result: dto.TimeEntryImpl{ID: "1"},
		},
		&simpleTestCase{
			name: "error response",
			param: api.CreateTimeEntryParam{
				Workspace: exampleID,
				Start: MustParseTime(timehlp.SimplerTimeFormat,
					"2022-11-07 10:00"),
			},

			requestMethod: "post",
			requestUrl:    uri,
			requestBody:   `{"start":"2022-11-07T10:00:00Z"}`,

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: `error`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.CreateTimeEntry(
					p.(api.CreateTimeEntryParam))
			})
	}
}
