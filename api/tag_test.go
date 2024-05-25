package api_test

import (
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

func TestGetTags(t *testing.T) {
	errPrefix := `get tags.*: `
	uri := "/v1/workspaces/" + exampleID +
		"/tags"
	var l []dto.Tag

	tts := []testCase{
		&simpleTestCase{
			name:  "requires workspace",
			param: api.GetTagsParam{},
			err:   errPrefix + "workspace is required",
		},
		&simpleTestCase{
			name:  "valid workspace",
			param: api.GetTagsParam{Workspace: "w"},
			err:   errPrefix + "workspace .* is not valid ID",
		},
		(&multiRequestTestCase{
			name: "get all pages, but find none",
			param: api.GetTagsParam{
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
			param: api.GetTagsParam{
				Workspace: exampleID,
				PaginationParam: api.PaginationParam{
					PageSize: 2,
					AllPages: true,
				},
			},

			result: []dto.Tag{
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
			param: api.GetTagsParam{
				Workspace:       exampleID,
				Name:            "tag",
				PaginationParam: api.AllPages(),
			},

			result: []dto.Tag{{ID: "p1", Name: "tag 1"}},

			requestMethod: "get",
			requestUrl: uri +
				"?name=tag&page=1&page-size=50",

			responseStatus: 200,
			responseBody:   `[{"id":"p1", "name": "tag 1"}]`,
		},
		&simpleTestCase{
			name: "error response",
			param: api.GetTagsParam{
				Workspace:       exampleID,
				PaginationParam: api.PaginationParam{Page: 2},
			},

			requestMethod: "get",
			requestUrl:    uri + "?page=2&page-size=50",

			responseStatus: 400,
			responseBody:   `{"code": 10, "message":"error"}`,

			err: errPrefix + `error \(code: 10\)`,
		},
	}

	for _, tt := range tts {
		runClient(t, tt,
			func(c api.Client, p interface{}) (interface{}, error) {
				return c.GetTags(
					p.(api.GetTagsParam))
			})
	}

}
