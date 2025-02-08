package router_test

import (
	"net/http"
	"testing"

	"github.com/cerfical/muzik/internal/httpserv/router"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTest))
}

type RouterTest struct {
	suite.Suite

	handler *mocks.Handler
	expect  *httpexpect.Expect
	router  *router.Router
}

func (t *RouterTest) SetupSubTest() {
	t.handler = mocks.NewHandler(t.T())
	t.router = router.New().
		Routes("/users/{userId}", []router.Endpoint{
			{"GET", t.handler.ServeHTTP},
		}).
		Routes("/users/", []router.Endpoint{
			{"GET", t.handler.ServeHTTP},
		}).
		Routes("/articles/{articleId}/comments/{commentId}", []router.Endpoint{
			{"GET", t.handler.ServeHTTP},
		}).
		Routes("/articles/", []router.Endpoint{
			{"PUT", t.handler.ServeHTTP},
			{"PATCH", t.handler.ServeHTTP},
		})

	t.expect = httpexpect.WithConfig(httpexpect.Config{
		TestName: t.T().Name(),
		Reporter: httpexpect.NewAssertReporter(t.T()),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(t.router),
		},
	})
}

func (t *RouterTest) TestRoutes() {
	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"200_ok", "/users/", http.StatusOK},
		{"404_nonexistent_path", "/user/", http.StatusNotFound},
		{"404_nonexistent_subpath", "/users/articles/", http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.status == http.StatusOK {
				t.handler.EXPECT().
					ServeHTTP(mock.Anything, mock.Anything).
					Return()
			}

			e := t.expect.GET(test.path).
				Expect()
			e.Status(test.status)
		})
	}
}

func (t *RouterTest) TestRoutesWithPatterns() {
	tests := []struct {
		name   string
		path   string
		values map[string]string
	}{
		{"single_param", "/users/9", map[string]string{
			"userId": "9",
		}},

		{"multiple_params", "/articles/9/comments/2", map[string]string{
			"articleId": "9",
			"commentId": "2",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			t.handler.EXPECT().
				ServeHTTP(mock.Anything, mock.Anything).
				Return().
				Run(func(_ http.ResponseWriter, r *http.Request) {
					for k, v := range test.values {
						t.Equal(v, r.PathValue(k))
					}
				})

			e := t.expect.GET(test.path).
				Expect()
			e.Status(http.StatusOK)
		})
	}
}

func (t *RouterTest) TestAllowedMethods() {
	tests := []struct {
		name   string
		method string
		path   string
		status int
		allow  string
	}{
		{"200_ok", "GET", "/users/", http.StatusOK, ""},
		{"405_unknown_method", "PUT", "/users/", http.StatusMethodNotAllowed, "GET, HEAD"},
		{"405_unknown_methods", "POST", "/articles/", http.StatusMethodNotAllowed, "PATCH, PUT"},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.status == http.StatusOK {
				t.handler.EXPECT().
					ServeHTTP(mock.Anything, mock.Anything).
					Return()
			}

			e := t.expect.Request(test.method, test.path).
				Expect()

			e.Status(test.status).
				Header("Allow").IsEqual(test.allow)
		})
	}
}
