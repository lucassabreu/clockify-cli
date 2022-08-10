package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	name   string
	param  interface{}
	result interface{}
	err    string

	requestMethod string
	requestUrl    string
	requestBody   string

	responseStatus int
	responseBody   string
}

var exampleID = "62f2af744a912b05acc7c79e"

func TestUpdateProject(t *testing.T) {
	bt := true
	bf := false
	n := "special"
	empty := ""
	tts := []testcase{
		{
			name:  "project is required",
			param: api.UpdateProjectParam{Workspace: "w"},
			err:   "update project: project id is required",
		},
		{
			name:  "workspace is required",
			param: api.UpdateProjectParam{ID: "p-1"},
			err:   "update project: workspace is required",
		},
		{
			name: "project should be a ID",
			param: api.UpdateProjectParam{
				ID:        "p-1",
				Workspace: exampleID,
			},
			err: "update project: project id (.*) is not valid",
		},
		{
			name: "workspace should be a ID",
			param: api.UpdateProjectParam{
				ID:        exampleID,
				Workspace: "w",
			},
			err: "update project: workspace (.*) is not valid",
		},
		{
			name: "color is not hex",
			param: api.UpdateProjectParam{
				ID:        exampleID,
				Workspace: exampleID,
				Color:     "#zzz",
			},
			err: "update project: color .* is not a hex string",
		},
		{
			name: "color must have 3 or 6 numbers (4)",
			param: api.UpdateProjectParam{
				ID:        exampleID,
				Workspace: exampleID,
				Color:     "#0000",
			},
			err: "update project: color must have 3.*or 6.*numbers",
		},
		{
			name: "color must have 3 or 6 numbers (2)",
			param: api.UpdateProjectParam{
				ID:        exampleID,
				Workspace: exampleID,
				Color:     "#00",
			},
			err: "update project: color must have 3.*or 6.*numbers",
		},
		{
			name: "empty update",
			param: api.UpdateProjectParam{
				ID:        exampleID,
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
				ID:        exampleID,
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
				ID:        exampleID,
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
				ID:        exampleID,
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
				ID:        exampleID,
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
				ID:        exampleID,
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
				ID:        exampleID,
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

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.requestUrl, r.URL.String())
				assert.Equal(t, tt.requestMethod, strings.ToLower(r.Method))

				var eMap, aMap map[string]interface{}
				assert.NoError(t,
					json.NewDecoder(r.Body).Decode(&aMap))
				assert.NoError(t,
					json.Unmarshal([]byte(tt.requestBody), &eMap))

				assert.Equal(t, eMap, aMap)

				w.WriteHeader(tt.responseStatus)
				if tt.responseBody == "" {
					tt.responseBody = "{}"
				}
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer s.Close()

			c, _ := api.NewClientFromUrlAndKey(
				"a-key",
				s.URL,
			)

			p, err := c.UpdateProject(tt.param.(api.UpdateProjectParam))
			if tt.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tt.err, err.Error())
				return
			}

			assert.NoError(t, err)
			if tt.result == nil {
				return
			}
			assert.Equal(t, tt.result, p)
		})
	}
}
