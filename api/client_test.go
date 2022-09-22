package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/stretchr/testify/assert"
)

var exampleID = "62f2af744a912b05acc7c79e"

func TestUpdateProjectMemberships(t *testing.T) {
	exampleID2 := "62f2af744a912b05acc7c792"
	errPrefix := `update project memberships: `
	uri := "/v1/workspaces/" + exampleID +
		"/projects/" + exampleID +
		"/memberships"

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.UpdateProjectMembershipsParam{ProjectID: "p1"},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "requires project",
			param: api.UpdateProjectMembershipsParam{Workspace: "w"},
			err:   errPrefix + "project id is required",
		},
		&simpleTestCase{
			name: "valid workspace",
			param: api.UpdateProjectMembershipsParam{
				Workspace: "w",
				ProjectID: exampleID,
			},
			err: errPrefix + "workspace .* is not valid ID",
		},
		&simpleTestCase{
			name: "valid project",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: "p1",
			},
			err: errPrefix + "project .* is not valid ID",
		},
		&simpleTestCase{
			name: "valid user or groups",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Memberships: []api.UpdateMembership{{
					UserOrGroupID: "ug",
				}},
			},
			err: errPrefix + "user or group .* is not valid ID",
		},
		&simpleTestCase{
			name: "required user or groups",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Memberships: []api.UpdateMembership{
					{UserOrGroupID: ""},
				},
			},
			err: errPrefix + `user or group is required`,
		},
		&simpleTestCase{
			name: "valid user or groups (second one)",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Memberships: []api.UpdateMembership{
					{UserOrGroupID: exampleID},
					{UserOrGroupID: "ug"},
				},
			},
			err: errPrefix + `user or group \("ug"\) is not valid ID`,
		},
		&simpleTestCase{
			name: "simplest update",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			result: dto.Project{ID: "p1", Name: "project 1"},

			requestMethod: "patch",
			requestUrl:    uri,
			requestBody:   `{"memberships":[]}`,

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "update with members",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Memberships: []api.UpdateMembership{{
					UserOrGroupID: exampleID,
				}},
			},

			result: dto.Project{ID: "p1", Name: "project 1"},

			requestMethod: "patch",
			requestUrl:    uri,
			requestBody: `{"memberships":[{
				"userId":"` + exampleID + `", "hourlyRate":{"amount":0}
			}]}`,

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "update with many members",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Memberships: []api.UpdateMembership{
					{UserOrGroupID: exampleID},
					{UserOrGroupID: exampleID2, HourlyRateAmount: 10},
				},
			},

			result: dto.Project{ID: "p1", Name: "project 1"},

			requestMethod: "patch",
			requestUrl:    uri,
			requestBody: `{"memberships":[
				{"userId":"` + exampleID + `", "hourlyRate":{"amount":0}},
				{"userId":"` + exampleID2 + `", "hourlyRate":{"amount":10}}
			]}`,

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.UpdateProjectMembershipsParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			requestMethod: "patch",
			requestUrl:    uri,
			requestBody:   `{"memberships":[]}`,

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: errPrefix + `error`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.UpdateProjectMemberships(
					p.(api.UpdateProjectMembershipsParam))
			})
	}
}

func TestDeleteProject(t *testing.T) {
	errPrefix := `delete project: `
	uri := "/v1/workspaces/" + exampleID + "/projects/" + exampleID

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.DeleteProjectParam{ProjectID: "p1"},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "requires project",
			param: api.DeleteProjectParam{Workspace: "w"},
			err:   errPrefix + "project id is required",
		},
		&simpleTestCase{
			name: "valid workspace",
			param: api.DeleteProjectParam{
				Workspace: "w",
				ProjectID: exampleID,
			},
			err: errPrefix + "workspace .* is not valid ID",
		},
		&simpleTestCase{
			name: "valid project",
			param: api.DeleteProjectParam{
				Workspace: exampleID,
				ProjectID: "p1",
			},
			err: errPrefix + "project .* is not valid ID",
		},
		&simpleTestCase{
			name: "delete",
			param: api.DeleteProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			result: dto.Project{ID: "p1", Name: "project 1"},

			requestMethod: "delete",
			requestUrl:    uri,

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.DeleteProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			requestMethod: "delete",
			requestUrl:    uri,

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: errPrefix + `error`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.DeleteProject(p.(api.DeleteProjectParam))
			})
	}
}

func TestGetProject(t *testing.T) {
	errPrefix := `get project "\w*": `
	uri := "/v1/workspaces/" + exampleID + "/projects/" + exampleID

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.GetProjectParam{ProjectID: "p1"},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "requires project",
			param: api.GetProjectParam{Workspace: "w"},
			err:   errPrefix + "project id is required",
		},
		&simpleTestCase{
			name: "valid workspace",
			param: api.GetProjectParam{
				Workspace: "w",
				ProjectID: exampleID,
			},
			err: errPrefix + "workspace .* is not valid ID",
		},
		&simpleTestCase{
			name: "valid project",
			param: api.GetProjectParam{
				Workspace: exampleID,
				ProjectID: "p1",
			},
			err: errPrefix + "project .* is not valid ID",
		},
		&simpleTestCase{
			name: "simple",
			param: api.GetProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			result: &dto.Project{ID: "p1", Name: "project 1"},

			requestMethod: "get",
			requestUrl:    uri,

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "hydrated",
			param: api.GetProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
				Hydrate:   true,
			},

			result: &dto.Project{ID: "p1", Name: "project 1", Hydrated: true},

			requestMethod: "get",
			requestUrl:    uri + "?hydrated=true",

			responseStatus: 200,
			responseBody:   `{"id":"p1", "name": "project 1"}`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.GetProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			requestMethod: "get",
			requestUrl:    uri,

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: errPrefix + `error`,
		},
		&simpleTestCase{
			name: "not found",
			param: api.GetProjectParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			requestMethod: "get",
			requestUrl:    uri,

			responseStatus: 404,
			responseBody:   `{"code": 0, "message":"not found"}`,

			err: errPrefix + `not found`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.GetProject(p.(api.GetProjectParam))
			})
	}
}

func TestGetProjects(t *testing.T) {
	errPrefix := "get projects: "
	uri := "/v1/workspaces/" + exampleID + "/projects"
	var l []dto.Project

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.GetProjectsParam{},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "valid workspace",
			param: api.GetProjectsParam{Workspace: "w"},
			err:   errPrefix + "workspace .* is not valid ID",
		},
		(&multiRequestTestCase{
			name: "get all pages, but find none",
			param: api.GetProjectsParam{
				Workspace:       exampleID,
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
			param: api.GetProjectsParam{
				Workspace: exampleID,
				PaginationParam: api.PaginationParam{
					PageSize: 2,
					AllPages: true,
				},
			},

			result: []dto.Project{
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
		(&multiRequestTestCase{
			name: "get all pages, hydrated",
			param: api.GetProjectsParam{
				Workspace: exampleID,
				Hydrate:   true,
				PaginationParam: api.PaginationParam{
					PageSize: 1,
					AllPages: true,
				},
			},

			result: []dto.Project{
				{ID: "p1", Hydrated: true},
				{ID: "p2", Hydrated: true},
			},
		}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?hydrated=true&page=1&page-size=1",
				status:   200,
				response: `[{"id":"p1"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?hydrated=true&page=2&page-size=1",
				status:   200,
				response: `[{"id":"p2"}]`,
			}).
			addHttpCall(&httpRequest{
				method:   "get",
				url:      uri + "?hydrated=true&page=3&page-size=1",
				status:   200,
				response: `[]`,
			}),
		&simpleTestCase{
			name: "all parameters",
			param: api.GetProjectsParam{
				Workspace:       exampleID,
				Hydrate:         true,
				Name:            "project",
				Clients:         []string{"c1", "c2"},
				PaginationParam: api.AllPages(),
			},

			result: []dto.Project{{
				ID: "p1", Name: "project 1", Hydrated: true}},

			requestMethod: "get",
			requestUrl: uri +
				"?clients=c1%2Cc2&hydrated=true&name=project&" +
				"page=1&page-size=50",

			responseStatus: 200,
			responseBody:   `[{"id":"p1", "name": "project 1"}]`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.GetProjectsParam{
				Workspace:       exampleID,
				PaginationParam: api.PaginationParam{Page: 2},
			},

			requestMethod: "get",
			requestUrl:    uri + "?page=2&page-size=50",

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: `get projects: error \(code: 10\)`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.GetProjects(p.(api.GetProjectsParam))
			})
	}
}

func TestUpdateProjectTemplate(t *testing.T) {
	errPrefix := "update project template: "
	tts := []simpleTestCase{
		{
			name: "workspace require",
			param: api.UpdateProjectTemplateParam{
				ProjectID: exampleID,
			},

			err: errPrefix + "workspace is required",
		},
		{
			name: "project require",
			param: api.UpdateProjectTemplateParam{
				Workspace: exampleID,
			},

			err: errPrefix + "project id is required",
		},
		{
			name: "valid workspace",
			param: api.UpdateProjectTemplateParam{
				ProjectID: exampleID,
				Workspace: "w",
			},

			err: errPrefix + "workspace .* is not valid ID",
		},
		{
			name: "valid project",
			param: api.UpdateProjectTemplateParam{
				ProjectID: "p",
				Workspace: exampleID,
			},

			err: errPrefix + "project .* is not valid ID",
		},
		{
			name: "into template",
			param: api.UpdateProjectTemplateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Template:  true,
			},

			result: dto.Project{ID: exampleID},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/template",
			requestBody: `{"isTemplate":true}`,

			responseStatus: 200,
			responseBody:   `{"id":"` + exampleID + `"}`,
		},
		{
			name: "not a template",
			param: api.UpdateProjectTemplateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Template:  false,
			},

			result: dto.Project{ID: exampleID},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/template",
			requestBody: `{"isTemplate":false}`,

			responseStatus: 200,
			responseBody:   `{"id":"` + exampleID + `"}`,
		},
		{
			name: "error",
			param: api.UpdateProjectTemplateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Template:  false,
			},

			err: errPrefix + "failed .code: 90.",

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/template",
			requestBody: `{"isTemplate":false}`,

			responseStatus: 400,
			responseBody:   `{"message":"failed", "code": 90}`,
		},
	}

	for i := range tts {
		runClient(t, &tts[i],
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.UpdateProjectTemplate(
					p.(api.UpdateProjectTemplateParam))
			})
	}
}

func TestUpdateProjectEstimate(t *testing.T) {
	errPrefix := "update project estimate: "
	tts := []simpleTestCase{
		{
			name: "workspace require",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Method:    api.EstimateMethodNone,
			},

			err: errPrefix + "workspace is required",
		},
		{
			name: "project require",
			param: api.UpdateProjectEstimateParam{
				Workspace: exampleID,
				Method:    api.EstimateMethodNone,
			},

			err: errPrefix + "project id is required",
		},
		{
			name: "estimate method required",
			param: api.UpdateProjectEstimateParam{
				Workspace: exampleID,
				ProjectID: exampleID,
			},

			err: errPrefix + "estimate method is required",
		},
		{
			name: "valid workspace",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: "w",
				Method:    api.EstimateMethodNone,
			},

			err: errPrefix + "workspace .* is not valid ID",
		},
		{
			name: "valid project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: "p",
				Workspace: exampleID,
				Method:    api.EstimateMethodNone,
			},

			err: errPrefix + "project .* is not valid ID",
		},
		{
			name: "valid method",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    "m",
			},

			err: errPrefix + "valid options for estimate method are",
		},
		{
			name: "type should be set for budget",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodBudget,
				Type:      "t",
			},

			err: errPrefix + "valid options for estimate type are",
		},
		{
			name: "valid reset option",
			param: api.UpdateProjectEstimateParam{
				ProjectID:   exampleID,
				Workspace:   exampleID,
				Method:      api.EstimateMethodBudget,
				Type:        api.EstimateTypeTask,
				ResetOption: "daily",
			},

			err: errPrefix + "valid options for reset option are",
		},
		{
			name: "type should be set for time",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
			},

			err: errPrefix + "valid options for estimate type are",
		},
		{
			name: "estimate should be set for budget method & type project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodBudget,
				Type:      api.EstimateTypeProject,
			},

			err: errPrefix +
				"estimate should be greater than zero for type project",
		},
		{
			name: "estimate should be set for time method & type project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
				Type:      api.EstimateTypeProject,
			},

			err: errPrefix +
				"estimate should be greater than zero for type project",
		},
		{
			name: "estimate should be positive for time method & type project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
				Type:      api.EstimateTypeProject,
				Estimate:  -1,
			},

			err: errPrefix +
				"estimate should be greater than zero for type project",
		},
		{
			name: "estimate should be positive for time method & type project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
				Type:      api.EstimateTypeProject,
				Estimate:  -1,
			},

			err: errPrefix +
				"estimate should be greater than zero for type project",
		},
		{
			name: "set estimate with budget for project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodBudget,
				Type:      api.EstimateTypeProject,
				Estimate:  1000,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"timeEstimate": {"active": false},
				"budgetEstimate": {
					"active": true,
					"estimate": 1000,
					"type": "MANUAL"
				}
			}`,

			responseStatus: 200,
		},
		{
			name: "set estimate with time for project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
				Type:      api.EstimateTypeProject,
				Estimate:  int64(time.Minute)*90 + int64(time.Second)*15,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"budgetEstimate": {"active": false},
				"timeEstimate": {
					"active": true,
					"estimate": "PT1H30M15S",
					"type": "MANUAL"
				}
			}`,

			responseStatus: 200,
		},
		{
			name: "set estimate to none for project",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodNone,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"budgetEstimate": {"active": false},
				"timeEstimate": {"active": false}
			}`,

			responseStatus: 200,
		},
		{
			name: "set estimate with budget for tasks",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodBudget,
				Type:      api.EstimateTypeTask,
				Estimate:  1000,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"timeEstimate": {"active": false},
				"budgetEstimate": {
					"active": true,
					"type": "AUTO"
				}
			}`,

			responseStatus: 200,
		},
		{
			name: "set estimate with time for task",
			param: api.UpdateProjectEstimateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Method:    api.EstimateMethodTime,
				Type:      api.EstimateTypeTask,
				Estimate:  int64(time.Minute)*90 + int64(time.Second)*15,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"budgetEstimate": {"active": false},
				"timeEstimate": {
					"active": true,
					"type": "AUTO"
				}
			}`,

			responseStatus: 200,
		},
		{
			name: "set estimate with time for task, and monthly reset",
			param: api.UpdateProjectEstimateParam{
				ProjectID:   exampleID,
				Workspace:   exampleID,
				Method:      api.EstimateMethodTime,
				Type:        api.EstimateTypeTask,
				ResetOption: api.EstimateResetOptionMonthly,
			},

			requestMethod: "patch",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID + "/estimate",
			requestBody: `{
				"budgetEstimate": {"active": false},
				"timeEstimate": {
					"active": true,
					"type": "AUTO",
					"resetOption": "MONTHLY"
				}
			}`,

			responseStatus: 200,
		},
	}

	for i := range tts {
		runClient(t, &tts[i],
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.UpdateProjectEstimate(
					p.(api.UpdateProjectEstimateParam))
			})
	}
}

func TestUpdateProjectUserCostRate(t *testing.T) {
	testUpdateProjectUserRate(t,
		"update project user cost rate: ",
		"cost-rate",
		func(c api.Client, p interface{}) (interface{}, error) {
			return c.UpdateProjectUserCostRate(
				p.(api.UpdateProjectUserRateParam))
		})
}

func TestUpdateProjectUserBillableRate(t *testing.T) {
	testUpdateProjectUserRate(t,
		"update project user billable rate: ",
		"hourly-rate",
		func(c api.Client, p interface{}) (interface{}, error) {
			return c.UpdateProjectUserBillableRate(
				p.(api.UpdateProjectUserRateParam))
		})
}

func testUpdateProjectUserRate(t *testing.T,
	errPrefix, uriSufix string,
	fn func(api.Client, interface{}) (interface{}, error)) {
	since, _ := time.Parse("2006-01-02", "2022-02-02")
	tts := []simpleTestCase{
		{
			name: "project is required",
			param: api.UpdateProjectUserRateParam{
				Workspace: "w",
				UserID:    "u",
			},
			err: errPrefix + "project id is required",
		},
		{
			name: "workspace is required",
			param: api.UpdateProjectUserRateParam{
				ProjectID: "p-1",
				UserID:    "u",
			},
			err: errPrefix + "workspace is required",
		},
		{
			name: "user is required",
			param: api.UpdateProjectUserRateParam{
				ProjectID: "p-1",
				Workspace: "w",
			},
			err: errPrefix + "user id is required",
		},
		{
			name: "project should be a ID",
			param: api.UpdateProjectUserRateParam{
				ProjectID: "p-1",
				Workspace: exampleID,
				UserID:    exampleID,
			},
			err: errPrefix + "project id (.*) is not valid",
		},
		{
			name: "user should be a ID",
			param: api.UpdateProjectUserRateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				UserID:    "u-1",
			},
			err: errPrefix + "user id (.*) is not valid",
		},
		{
			name: "workspace should be a ID",
			param: api.UpdateProjectUserRateParam{
				ProjectID: exampleID,
				Workspace: "w",
				UserID:    exampleID,
			},
			err: errPrefix + "workspace (.*) is not valid",
		},
		{
			name: "only amount",
			param: api.UpdateProjectUserRateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				UserID:    exampleID,
				Amount:    10,
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID +
				"/users/" + exampleID + "/" + uriSufix,
			requestBody: `{"amount":10}`,

			responseStatus: 200,
		},
		{
			name: "amount and since",
			param: api.UpdateProjectUserRateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				UserID:    exampleID,
				Amount:    10,
				Since:     &since,
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID +
				"/users/" + exampleID + "/" + uriSufix,
			requestBody: `{"amount":10,"since":"2022-02-02T00:00:00Z"}`,

			err:            errPrefix + "custom error.*code: 42",
			responseStatus: 400,
			responseBody:   `{"message":"custom error","code":42}`,
		},
		{
			name: "fail",
			param: api.UpdateProjectUserRateParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				UserID:    exampleID,
				Amount:    10,
				Since:     &since,
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID +
				"/users/" + exampleID + "/" + uriSufix,
			requestBody: `{"amount":10,"since":"2022-02-02T00:00:00Z"}`,

			err:            errPrefix + "custom error.*code: 42",
			responseStatus: 400,
			responseBody:   `{"message":"custom error","code":42}`,
		},
	}

	for i := range tts {
		runClient(t, &tts[i], fn)
	}
}

func TestUpdateProject(t *testing.T) {
	bt := true
	bf := false
	n := "special"
	empty := ""
	tts := []simpleTestCase{
		{
			name:  "project is required",
			param: api.UpdateProjectParam{Workspace: "w"},
			err:   "update project: project id is required",
		},
		{
			name:  "workspace is required",
			param: api.UpdateProjectParam{ProjectID: "p-1"},
			err:   "update project: workspace is required",
		},
		{
			name: "project should be a ID",
			param: api.UpdateProjectParam{
				ProjectID: "p-1",
				Workspace: exampleID,
			},
			err: "update project: project id (.*) is not valid",
		},
		{
			name: "workspace should be a ID",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: "w",
			},
			err: "update project: workspace (.*) is not valid",
		},
		{
			name: "color is not hex",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Color:     "#zzz",
			},
			err: "update project: color .* is not a hex string",
		},
		{
			name: "color must have 3 or 6 numbers (4)",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Color:     "#0000",
			},
			err: "update project: color must have 3.*or 6.*numbers",
		},
		{
			name: "color must have 3 or 6 numbers (2)",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Color:     "#00",
			},
			err: "update project: color must have 3.*or 6.*numbers",
		},
		{
			name: "empty update",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: "{}",

			responseStatus: 200,
		},
		{
			name: "full update",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				Name:      "a new name",
				Public:    &bt,
				Archived:  &bf,
				Note:      &n,
				ClientId:  &exampleID,
				Color:     "012345",
				Billable:  &bt,
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{
				"archived":false,
				"isPublic":true,
				"billable":true,
				"clientId":"` + exampleID + `",
				"note": "special",
				"color": "#012345",
				"name":"a new name"
			}`,

			responseStatus: 200,
		},
		{
			name: "expand color and remove client",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
				ClientId:  &empty,
				Color:     "#0f0",
			},

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{
				"clientId":"",
				"color": "#00ff00"
			}`,

			responseStatus: 200,
		},
		{
			name: "report 404",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
			},
			err: "update project: Nothing was found .*404",

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{}`,

			responseStatus: 404,
		},
		{
			name: "report 403",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
			},
			err: "update project: Forbidden.*403",

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{}`,

			responseStatus: 403,
		},
		{
			name: "report no response",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
			},
			err: "update project: No response",

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{}`,

			responseStatus: 400,
			responseBody:   `{}`,
		},
		{
			name: "report error",
			param: api.UpdateProjectParam{
				ProjectID: exampleID,
				Workspace: exampleID,
			},
			err: "update project: custom error.*code: 42",

			requestMethod: "put",
			requestUrl: "/v1/workspaces/" + exampleID +
				"/projects/" + exampleID,
			requestBody: `{}`,

			responseStatus: 400,
			responseBody:   `{"message":"custom error","code":42}`,
		},
	}

	for i := range tts {
		runClient(t, &tts[i], func(
			c api.Client, p interface{}) (interface{}, error) {
			return c.UpdateProject(p.(api.UpdateProjectParam))
		})
	}
}

type testCase interface {
	getName() string
	getParam() interface{}
	getResult() interface{}
	getErr() string

	hasHttpCalls() bool
	getHttpCallFor(uri string) httpCall
	getPendingHttpCalls() []httpCall
}

type httpCall interface {
	getRequestMethod() string
	getRequestUrl() string
	getRequestBody() string
	getResponseStatus() int
	getResponseBody() string
}

func runClient(t *testing.T, tt testCase,
	fn func(api.Client, interface{}) (interface{}, error)) {

	t.Run(tt.getName(), func(t *testing.T) {
		httpCalled := false
		t.Cleanup(func() {
			if !tt.hasHttpCalls() {
				assert.False(t, httpCalled, "should not call api")
				return
			}
			assert.True(t, httpCalled, "should call api")
		})
		s := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				httpCalled = true
				if !tt.hasHttpCalls() {
					t.Error("should not call api")
					w.WriteHeader(500)
					return
				}

				hc := tt.getHttpCallFor(r.URL.String())
				if hc == nil {
					assert.FailNow(t, "should not call api "+r.URL.String())
					w.WriteHeader(500)
					return
				}

				assert.Equal(t, hc.getRequestUrl(), r.URL.String())
				assert.Equal(t,
					hc.getRequestMethod(), strings.ToLower(r.Method))

				b, _ := io.ReadAll(r.Body)
				if hc.getRequestBody() != "" {
					var eMap, aMap map[string]interface{}
					assert.NoError(t, json.Unmarshal(b, &aMap))
					assert.NoError(t,
						json.Unmarshal([]byte(hc.getRequestBody()), &eMap))

					assert.Equal(t, eMap, aMap)
				} else {
					assert.Empty(t, string(b))
				}

				w.WriteHeader(hc.getResponseStatus())
				rb := hc.getResponseBody()
				if rb == "" {
					rb = "{}"
				}
				_, err := w.Write([]byte(rb))
				assert.NoError(t, err)
			}))
		defer s.Close()

		c, _ := api.NewClientFromUrlAndKey(
			"a-key",
			s.URL,
		)

		r, err := fn(c, tt.getParam())
		if tt.getErr() != "" {
			if !assert.Error(t, err) {
				return
			}
			assert.Regexp(t, tt.getErr(), err.Error())
			return
		}

		if !assert.NoError(t, err) || tt.getResult() == nil {
			return
		}
		assert.Equal(t, tt.getResult(), r)
	})
}

type simpleTestCase struct {
	name   string
	param  interface{}
	result interface{}
	err    string

	requestMethod string
	requestUrl    string
	requestBody   string

	responseStatus int
	responseBody   string

	once bool
}

func (s *simpleTestCase) getRequestMethod() string {
	return s.requestMethod
}

func (s *simpleTestCase) getRequestUrl() string {
	return s.requestUrl
}

func (s *simpleTestCase) getRequestBody() string {
	return s.requestBody
}

func (s *simpleTestCase) getResponseStatus() int {
	return s.responseStatus
}

func (s *simpleTestCase) getResponseBody() string {
	return s.responseBody
}

func (s *simpleTestCase) getName() string {
	return s.name
}

func (s *simpleTestCase) getParam() interface{} {
	return s.param
}

func (s *simpleTestCase) getResult() interface{} {
	return s.result
}

func (s *simpleTestCase) getErr() string {
	return s.err
}

func (s *simpleTestCase) getHttpCallFor(_ string) httpCall {
	if !s.once {
		s.once = true
		return s
	}
	return nil
}

func (s *simpleTestCase) getPendingHttpCalls() []httpCall {
	if s.once {
		return []httpCall{}
	}

	return []httpCall{s}
}

func (s *simpleTestCase) hasHttpCalls() bool {
	return s.requestUrl != ""
}

type multiRequestTestCase struct {
	name  string
	param interface{}

	err    string
	result interface{}

	calls    map[string]httpCall
	hasCalls bool
}

func (m *multiRequestTestCase) getName() string {
	return m.name
}

func (m *multiRequestTestCase) getParam() interface{} {
	return m.param
}

func (m *multiRequestTestCase) getResult() interface{} {
	return m.result
}

func (m *multiRequestTestCase) getErr() string {
	return m.err
}

func (m *multiRequestTestCase) hasHttpCalls() bool {
	return m.hasCalls
}

func (m *multiRequestTestCase) getHttpCallFor(uri string) httpCall {
	if !m.hasCalls {
		return nil
	}
	c := m.calls[uri]
	delete(m.calls, uri)
	return c
}

func (m *multiRequestTestCase) getPendingHttpCalls() []httpCall {
	if !m.hasCalls {
		return []httpCall{}
	}
	l := make([]httpCall, len(m.calls))
	for _, c := range m.calls {
		l = append(l, c)
	}
	return l
}

func (m *multiRequestTestCase) addHttpCall(c httpCall) *multiRequestTestCase {
	if m.calls == nil {
		m.calls = make(map[string]httpCall)
		m.hasCalls = true
	}

	if _, ok := m.calls[c.getRequestUrl()]; ok {
		panic("http call for " + c.getRequestUrl() + " already exists")
	}
	m.calls[c.getRequestUrl()] = c
	return m
}

type httpRequest struct {
	method string
	url    string
	body   string

	status   int
	response string
}

func (h *httpRequest) getRequestMethod() string {
	return h.method
}

func (h *httpRequest) getRequestUrl() string {
	return h.url
}

func (h *httpRequest) getRequestBody() string {
	return h.body
}

func (h *httpRequest) getResponseStatus() int {
	return h.status
}

func (h *httpRequest) getResponseBody() string {
	return h.response
}
