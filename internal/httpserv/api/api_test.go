package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITest))
}

type APITest struct {
	suite.Suite

	store  *mocks.TrackStore
	expect *httpexpect.Expect
}

func (t *APITest) SetupSubTest() {
	t.store = mocks.NewTrackStore(t.T())
	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		BaseURL:  "/api/tracks",
		Reporter: httpexpect.NewAssertReporter(t.T()),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api.NewHandler(t.store, nil)),
		},
	})
}

func (t *APITest) TestContentTypeHeaderCheck() {
	tests := []struct {
		name        string
		contentType string
	}{
		{"ok_application_json", "application/json"},
		{"ok_application_json_valid_charset", "application/json;charset=utf-8"},

		{"fail_invalid_charset", "application/json;charset=utf8"},
		{"fail_unknown_param", "application/json;q=1"},
		{"fail_invalid_param", "application/json;charset"},
		{"fail_unknown_type", "application/xml"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Content-Type", test.contentType)

			if strings.HasPrefix(test.name, "ok") {
				t.anyStatusExcept(e.Expect(), http.StatusUnsupportedMediaType)
			} else {
				e := e.WithMatcher(isErrorResponse(http.StatusUnsupportedMediaType)).
					Expect()

				e.Header("Accept-Post").
					IsEqual("application/json")

				e.JSON().
					Path("$.error.source.header").
					IsEqual("Content-Type")
			}
		})
	}
}

func (t *APITest) TestAcceptHeaderCheck() {
	tests := []struct {
		name   string
		accept string
	}{
		{"ok_application_json", "application/json"},
		{"ok_application_any", "application/*"},
		{"ok_any_any", "*/*"},
		{"ok_any_json", "*/json"},
		{"ok_param_valid_qvalue", "application/json;q=1"},

		{"fail_unknown_type", "application/xml"},
		{"fail_invalid_type", "application/"},
		{"fail_param_zero_qvalue", "application/json;q=0"},
		{"fail_param_invalid_qvalue", "application/json;q=str"},
		{"fail_param_invalid_qvalue", "application/json;q=str"},
		{"fail_unknown_param", "application/json;p=1"},
		{"fail_invalid_param", "application/json;p"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Accept", test.accept)

			if strings.HasPrefix(test.name, "ok") {
				t.anyStatusExcept(e.Expect(), http.StatusNotAcceptable)
			} else {
				e.WithMatcher(isErrorResponse(http.StatusNotAcceptable)).
					Expect().
					JSON().
					Path("$.error.source.header").
					IsEqual("Accept")
			}
		})
	}
}

func (t *APITest) TestAllowedMethods() {
	tests := []struct {
		name string

		method, path string
		allow        string
	}{
		{"ok_post", "POST", "/", ""},

		{"fail_post_by_id", "POST", "/1", "GET"},
		{"fail_put_all", "PUT", "/", "GET, POST"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.Request(test.method, test.path)
			if strings.HasPrefix(test.name, "ok") {
				t.anyStatusExcept(e.Expect(), http.StatusMethodNotAllowed)
			} else {
				e.Expect().
					Status(http.StatusMethodNotAllowed).
					Header("Allow").
					IsEqual(test.allow)
			}
		})
	}
}

func (t *APITest) anyStatusExcept(r *httpexpect.Response, status int) {
	t.NotEqual(status, r.Raw().StatusCode)
}

func isErrorResponse(status int) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		obj := r.Status(status).
			JSON().
			Schema(modelRef("error")).
			Object()

		obj.Path("$.error.status").
			IsEqual(strconv.Itoa(status))
	}
}

func isDataResponse(model string, status int, data any) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		r.Status(status).
			JSON().
			Schema(modelRef(model)).
			Object().
			ContainsValue(data)
	}
}

func modelRef(model string) string {
	p, err := filepath.Abs("../../../api/models.json")
	if err != nil {
		panic(err)
	}

	p, err = url.JoinPath("/", p)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("file://%s#/$def/%s", p, model)
}
