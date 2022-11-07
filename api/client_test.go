package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/stretchr/testify/assert"
)

var exampleID = "62f2af744a912b05acc7c79e"

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
