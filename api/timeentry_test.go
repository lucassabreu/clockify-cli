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

func TestGetUsersHydratedTimeEntries(t *testing.T) {
	otherID := "62f2af744a912b05acc7c79f"
	meUri := "/v1/user"
	entriesUri := func(user string) string {
		return "/v1/workspaces/" + exampleID + "/user/" + user +
			"/time-entries?hydrated=1&page=1&page-size=50"
	}
	usersUri := "/v1/workspaces/" + exampleID + "/users?page=1&page-size=50"

	me := dto.User{ID: exampleID, Name: "me"}
	other := dto.User{ID: otherID, Name: "other"}
	meBody := `{"id":"` + exampleID + `","name":"me"}`

	tts := []testCase{
		(&multiRequestTestCase{
			name: "hydrates with the token owner",
			param: api.GetUserTimeEntriesParam{
				Workspace:       exampleID,
				UserID:          exampleID,
				PaginationParam: api.AllPages(),
			},

			result: []dto.TimeEntry{{ID: "t1", User: &me}},
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      entriesUri(exampleID),
				status:   200,
				response: `[{"id":"t1"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      meUri,
				status:   200,
				response: meBody,
			}),
		(&multiRequestTestCase{
			name: "hydrates another user from the workspace",
			param: api.GetUserTimeEntriesParam{
				Workspace:       exampleID,
				UserID:          otherID,
				PaginationParam: api.AllPages(),
			},

			result: []dto.TimeEntry{{ID: "t1", User: &other}},
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      entriesUri(otherID),
				status:   200,
				response: `[{"id":"t1"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      meUri,
				status:   200,
				response: meBody,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      usersUri,
				status:   200,
				response: `[{"id":"` + otherID + `","name":"other"}]`,
			}),
		(&multiRequestTestCase{
			name: "fails when the token owner is unknown",
			param: api.GetUserTimeEntriesParam{
				Workspace:       exampleID,
				UserID:          exampleID,
				PaginationParam: api.AllPages(),
			},

			err: `Access Denied \(code: 501\)`,
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      entriesUri(exampleID),
				status:   200,
				response: `[{"id":"t1"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      meUri,
				status:   403,
				response: `{"code": 501, "message":"Access Denied"}`,
			}),
		(&multiRequestTestCase{
			name: "fails when the workspace users can't be listed",
			param: api.GetUserTimeEntriesParam{
				Workspace:       exampleID,
				UserID:          otherID,
				PaginationParam: api.AllPages(),
			},

			err: `get users.*: Access Denied \(code: 501\)`,
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      entriesUri(otherID),
				status:   200,
				response: `[{"id":"t1"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      meUri,
				status:   200,
				response: meBody,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      usersUri,
				status:   403,
				response: `{"code": 501, "message":"Access Denied"}`,
			}),
		&simpleTestCase{
			name: "fails when the time entries can't be listed",
			param: api.GetUserTimeEntriesParam{
				Workspace:       exampleID,
				UserID:          exampleID,
				PaginationParam: api.AllPages(),
			},

			requestMethod: "get",
			requestUrl:    entriesUri(exampleID),

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: `get time entries from user .*: error \(code: 10\)`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.GetUsersHydratedTimeEntries(
					p.(api.GetUserTimeEntriesParam))
			})
	}
}
