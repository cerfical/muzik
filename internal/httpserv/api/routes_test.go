package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/cerfical/muzik/internal/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/mock"
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

	e := t.store.EXPECT()
	e.GetTracks(mock.Anything).
		Return([]model.Track{}, nil).
		Maybe()
	e.GetTrack(mock.Anything, mock.Anything).
		Return(&model.Track{}, nil).
		Maybe()
	e.DeleteTrack(mock.Anything, mock.Anything).
		Return(nil).
		Maybe()
	e.CreateTrack(mock.Anything, &model.TrackAttrs{}).
		Return(&model.Track{}, nil).
		Maybe()
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

			e.Status(http.StatusBadRequest)
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

			e.Status(http.StatusUnsupportedMediaType).
				Header("Accept-Post").IsEqual("application/json")
			e.JSON().Schema(errorResponse())
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
			e := t.expect.GET("/").
				WithHeader("Accept", test.accept).
				Expect()

			e.Status(http.StatusOK)
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
			e.JSON().Schema(errorResponse())
		})
	}
}

func trackDataResponse() string {
	return schema("TrackDataResponse")
}

func tracksDataResponse() string {
	return schema("TracksDataResponse")
}

func errorResponse() string {
	return schema("ErrorResponse")
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

	return fmt.Sprintf("file://%s#/$defs/%s", p, name)
}
