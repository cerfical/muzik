package api_test

import (
	"net/http"
	"testing"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/cerfical/muzik/internal/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var sampleTracks = []model.Track{
	{
		ID: 1,
		Attrs: model.TrackAttrs{
			Title: "Example Track #1",
		},
	},
	{
		ID: 2,
		Attrs: model.TrackAttrs{
			Title: "Example Track #2",
		},
	},
}

func TestTracks(t *testing.T) {
	suite.Run(t, new(TracksTest))
}

type TracksTest struct {
	suite.Suite

	store  *mocks.TrackStore
	expect *httpexpect.Expect
}

func (t *TracksTest) SetupTest() {
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

func (t *TracksTest) TestTracks_Get_Ok() {
	var response struct {
		Data model.Track `json:"data"`
	}
	response.Data = sampleTracks[0]

	t.store.EXPECT().
		TrackByID(mock.Anything, 1).
		Return(&sampleTracks[0], nil)

	t.expect.GET("/1").
		Expect().
		Status(http.StatusOK).
		JSON().
		Schema(trackDataResponseSchema()).
		IsEqual(&response)
}

func (t *TracksTest) TestTracks_Get_NotFound() {
	t.store.EXPECT().
		TrackByID(mock.Anything, 3).
		Return(nil, model.ErrNotFound)

	t.expect.GET("/3").
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Schema(errorSchema())
}

func (t *TracksTest) TestTracks_GetAll_Ok() {
	var response struct {
		Data []model.Track `json:"data"`
	}
	response.Data = sampleTracks

	t.store.EXPECT().
		AllTracks(mock.Anything).
		Return(sampleTracks, nil)

	t.expect.GET("/").
		Expect().
		Status(http.StatusOK).
		JSON().
		Schema(tracksDataResponseSchema()).
		IsEqual(&response)
}

func (t *TracksTest) TestTracks_Create_Created() {
	var request struct {
		Data struct {
			Attrs model.TrackAttrs `json:"attributes"`
		} `json:"data"`
	}
	request.Data.Attrs = sampleTracks[0].Attrs

	var response struct {
		Data struct {
			Attrs model.TrackAttrs `json:"attributes"`
			ID    int              `json:"id,string"`
		} `json:"data"`
	}
	response.Data.Attrs = sampleTracks[0].Attrs
	response.Data.ID = sampleTracks[0].ID

	t.store.EXPECT().
		CreateTrack(mock.Anything, &sampleTracks[0].Attrs).
		Return(&sampleTracks[0], nil)

	e := t.expect.POST("/").
		WithJSON(&request).
		Expect()

	e.Header("Location").IsEqual("/api/tracks/1")
	e.Status(http.StatusCreated).
		JSON().
		Schema(trackDataResponseSchema()).
		IsEqual(&response)
}

func (t *TracksTest) TestTracks_Create_BadRequest() {
	e := t.expect.POST("/").
		WithJSON("").
		Expect()

	e.Status(http.StatusBadRequest).
		JSON().
		Schema(errorSchema())
}

func trackDataResponseSchema() string {
	return schema("TrackResource")
}

func tracksDataResponseSchema() string {
	return schema("TracksResource")
}
