package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

func TestRoutes(t *testing.T) {
	suite.Run(t, new(RoutesTest))
}

type RoutesTest struct {
	suite.Suite

	store  *mocks.TrackStore
	expect *httpexpect.Expect
}

func (t *RoutesTest) SetupSubTest() {
	t.store = mocks.NewTrackStore(t.T())
	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		Reporter: httpexpect.NewAssertReporter(t.T()),
		BaseURL:  "/api/tracks/",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api.NewHandler(t.store, nil)),
		},
	})
}

func (t *RoutesTest) TestContentTypeCheck_Ok() {
	tests := []struct {
		name        string
		contentType string
	}{
		{"json", "application/json"},
		{"json_with_utf8", "application/json;charset=utf-8"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Content-Type", test.contentType).
				Expect()
			t.NotEqual(http.StatusUnsupportedMediaType, e.Raw().StatusCode)
		})
	}
}

func (t *RoutesTest) TestContentTypeCheck_Fail() {
	tests := []struct {
		name        string
		contentType string
	}{
		{"json_with_invalid_charset", "application/json;charset=utf8"},
		{"json_with_unknown_param", "application/json;q=1"},
		{"json_with_invalid_param", "application/json;q"},
		{"unknown_type", "application/xml"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Content-Type", test.contentType).
				Expect()

			e.Header("Accept-Post").IsEqual("application/json")
			e.Status(http.StatusUnsupportedMediaType).
				JSON().
				Schema(errorSchema())
		})
	}
}

func (t *RoutesTest) TestAcceptHeaderCheck_Ok() {
	tests := []struct {
		name   string
		accept string
	}{
		{"json", "application/json"},
		{"any_type", "*/*"},
		{"nonzero_qvalue", "application/json;q=1"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Accept", test.accept).
				Expect()
			t.NotEqual(http.StatusNotAcceptable, e.Raw().StatusCode)
		})
	}
}

func (t *RoutesTest) TestAcceptHeaderCheck_Fail() {
	tests := []struct {
		name   string
		accept string
	}{
		{"unknown_type", "application/xml"},
		{"invalid_type", "application/"},
		{"zero_qvalue", "application/json;q=str"},
		{"invalid_qvalue", "application/json;q=str"},
		{"unknown_param", "application/json;p=1"},
		{"invalid_param", "application/json;p"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.POST("/").
				WithHeader("Accept", test.accept).
				Expect()

			e.Status(http.StatusNotAcceptable)
			e.JSON().Schema(errorSchema())
		})
	}
}

func (t *RoutesTest) TestAllowedMethodsCheck_Ok() {
	tests := []struct {
		name         string
		method, path string
	}{
		{"post", "POST", "/"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.Request(test.method, test.path).
				Expect()
			t.NotEqual(http.StatusMethodNotAllowed, e.Raw().Status)
		})
	}
}

func (t *RoutesTest) TestAllowedMethodsCheck_Fail() {
	tests := []struct {
		name         string
		method, path string
		allow        string
	}{
		{"post_to_id", "POST", "/1", "GET"},
		{"put", "PUT", "/", "GET, POST"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			e := t.expect.Request(test.method, test.path).
				Expect()

			e.JSON().Schema(errorSchema())
			e.Status(http.StatusMethodNotAllowed).
				Header("Allow").IsEqual(test.allow)
		})
	}
}

func trackDataResponseSchema() string {
	return schema("TrackResource")
}

func tracksDataResponseSchema() string {
	return schema("TracksResource")
}

func errorSchema() string {
	return schema("Error")
}

func schema(name string) string {
	p, err := filepath.Abs("../../../api/models.json")
	if err != nil {
		panic(err)
	}

	p, err = url.JoinPath("/", p)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("file://%s#/$def/%s", p, name)
}
