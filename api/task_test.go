package api_test

import (
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

func TestGetTasks(t *testing.T) {
	errPrefix := `get tasks from project .*: `
	uri := "/v1/workspaces/" + exampleID +
		"/projects/" + exampleID +
		"/tasks"
	var l []dto.Task

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.GetTasksParam{ProjectID: exampleID},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "valid workspace",
			param: api.GetTasksParam{Workspace: "w", ProjectID: exampleID},
			err:   errPrefix + "workspace .* is not valid ID",
		},
		&simpleTestCase{
			name:  "requires project id",
			param: api.GetTasksParam{Workspace: exampleID},
			err:   errPrefix + "project id is required",
		},
		&simpleTestCase{
			name:  "valid project id",
			param: api.GetTasksParam{ProjectID: "w", Workspace: exampleID},
			err:   errPrefix + "project id .* is not valid ID",
		},
		(&multiRequestTestCase{
			name: "get all pages, but find none",
			param: api.GetTasksParam{
				Workspace:       exampleID,
				ProjectID:       exampleID,
				PaginationParam: api.AllPages(),
			},

			result: l,
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?page=1&page-size=50",
				status:   200,
				response: "[]",
			}),
		(&multiRequestTestCase{
			name: "get all pages, find five",
			param: api.GetTasksParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				PaginationParam: api.PaginationParam{
					PageSize: 2,
					AllPages: true,
				},
			},

			result: []dto.Task{
				{ID: "p1"},
				{ID: "p2"},
				{ID: "p3"},
				{ID: "p4"},
				{ID: "p5"},
			},
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?page=1&page-size=2",
				status:   200,
				response: `[{"id":"p1"},{"id":"p2"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?page=2&page-size=2",
				status:   200,
				response: `[{"id":"p3"},{"id":"p4"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?page=3&page-size=2",
				status:   200,
				response: `[{"id":"p5"}]`,
			}),
		&simpleTestCase{
			name: "all parameters",
			param: api.GetTasksParam{
				Workspace:       exampleID,
				ProjectID:       exampleID,
				Name:            "project",
				PaginationParam: api.AllPages(),
			},

			result: []dto.Task{{ID: "p1", Name: "project 1"}},

			requestMethod: "get",
			requestUrl: uri +
				"?name=project&page=1&page-size=50",

			responseStatus: 200,
			responseBody:   `[{"id":"p1", "name": "project 1"}]`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.GetTasksParam{
				Workspace:       exampleID,
				ProjectID:       exampleID,
				PaginationParam: api.PaginationParam{Page: 2},
			},

			requestMethod: "get",
			requestUrl:    uri + "?page=2&page-size=50",

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: `get tasks from project .*: error \(code: 10\)`,
		},
		&simpleTestCase{
			name: "missing estimate",
			param: api.GetTasksParam{
				Workspace:       exampleID,
				ProjectID:       exampleID,
				PaginationParam: api.AllPages(),
			},

			result: []dto.Task{
				{
					ID:           "wod",
					Name:         "without durations",
					ProjectID:    "p",
					Duration:     nil,
					Estimate:     nil,
					Status:       api.TaskStatusDone,
					UserGroupIDs: []string{},
					AssigneeIDs:  []string{},
					Billable:     true,
				},
				{
					ID:        "wd",
					Name:      "with durations",
					ProjectID: "p",
					Duration: &dto.Duration{
						Duration: 120 * time.Hour,
					},
					Estimate: &dto.Duration{
						Duration: 120 * time.Hour,
					},
					Status:       api.TaskStatusActive,
					UserGroupIDs: []string{},
					AssigneeIDs:  []string{},
					Billable:     true,
				},
			},

			requestMethod: "get",
			requestUrl:    uri + "?page=1&page-size=50",

			responseStatus: 200,
			responseBody: `[{
				"id": "wod",
				"name": "without durations",
				"projectId": "p",
				"assigneeIds": [],
				"assigneeId": null,
				"userGroupIds": [],
				"estimate": null,
				"status": "DONE",
				"duration": null,
				"billable": true,
				"hourlyRate": null,
				"costRate": null
			},{
				"id": "wd",
				"name": "with durations",
				"projectId": "p",
				"assigneeIds": [],
				"assigneeId": null,
				"userGroupIds": [],
				"estimate": "P120H",
				"status": "ACTIVE",
				"duration": "P120H",
				"billable": true,
				"hourlyRate": null,
				"costRate": null
			}]`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.GetTasks(
					p.(api.GetTasksParam))
			})
	}

}
