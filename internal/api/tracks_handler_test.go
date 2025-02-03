package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/cerfical/muzik/internal/mocks"
	"github.com/cerfical/muzik/internal/model"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

func TestTracks(t *testing.T) {
	suite.Run(t, new(TracksHandlerTest))
}

type TracksHandlerTest struct {
	suite.Suite

	store  *mocks.TrackStore
	expect *httpexpect.Expect
}

func (t *TracksHandlerTest) SetupSubTest() {
	t.store = mocks.NewTrackStore(t.T())
	tracks := tracksHandler{store: t.store}

	router := http.NewServeMux()
	router.HandleFunc("GET /{id}", tracks.get)
	router.HandleFunc("GET /", tracks.getAll)
	router.HandleFunc("POST /", tracks.create)

	client := http.Client{Transport: httpexpect.NewBinder(router)}
	t.expect = httpexpect.WithConfig(httpexpect.Config{
		Client:   &client,
		TestName: t.T().Name(),
		Reporter: httpexpect.NewAssertReporter(t.T()),
	})
}

func (t *TracksHandlerTest) TestTracksHandler_Get() {
	track := model.Track{ID: 7, Title: "Example Track"}
	tests := []struct {
		name      string
		id        any
		storeData *model.Track
		storeErr  error
		matcher   func(*httpexpect.Response)
	}{
		{"200_ok", track.ID, &track, nil, trackDataMatcher(http.StatusOK, &track)},
		{"404_not_found", track.ID, nil, model.ErrNotFound, errorMatcher(http.StatusNotFound)},
		{"404_invalid_id", "invalid_id", nil, nil, errorMatcher(http.StatusNotFound)},
		{"500_storage_fail", track.ID, nil, errors.New(""), errorMatcher(http.StatusInternalServerError)},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if id, ok := test.id.(int); ok {
				t.store.EXPECT().
					TrackByID(id).
					Return(test.storeData, test.storeErr)
			}

			t.expect.GET("/{id}", test.id).
				WithMatcher(test.matcher).
				Expect()
		})
	}
}

func (t *TracksHandlerTest) TestTracksHandler_GetAll() {
	tracks := []model.Track{
		{ID: 7, Title: "Example Track #1"},
		{ID: 2, Title: "Example Track #2"},
	}
	empty := make([]model.Track, 0)

	tests := []struct {
		name      string
		storeData []model.Track
		storeErr  error
		matcher   func(*httpexpect.Response)
	}{
		{"200_ok", tracks, nil, tracksDataMatcher(http.StatusOK, tracks)},
		{"200_empty", empty, nil, tracksDataMatcher(http.StatusOK, empty)},
		{"500_storage_fail", nil, errors.New(""), errorMatcher(http.StatusInternalServerError)},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			t.store.EXPECT().
				AllTracks().
				Return(test.storeData, test.storeErr)

			t.expect.GET("/").
				WithMatcher(test.matcher).
				Expect()
		})
	}
}

func (t *TracksHandlerTest) TestTracksHandler_Create() {
	track := model.Track{ID: 7, Title: "Example Track"}
	body := struct {
		Data model.Track `json:"data"`
	}{
		Data: track,
	}

	tests := []struct {
		name      string
		storeData *model.Track
		storeErr  error
		body      any
		location  string
		matcher   func(*httpexpect.Response)
	}{
		{"201_ok", &track, nil, body, "/7", trackDataMatcher(http.StatusCreated, &track)},
		{"400_bad_request", nil, nil, "{}", "", errorMatcher(http.StatusBadRequest)},
		{"500_storage_fail", &track, errors.New(""), body, "", errorMatcher(http.StatusInternalServerError)},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.storeData != nil {
				t.store.EXPECT().
					CreateTrack(test.storeData).
					Return(test.storeErr)
			}

			t.expect.POST("/").
				WithMatcher(test.matcher).
				WithJSON(test.body).
				Expect().
				Header("Location").
				IsEqual(test.location)
		})
	}
}

func errorMatcher(status int) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		obj := r.Status(status).
			JSON().
			Schema(modelRef("error")).
			Object()

		obj.Path("$.error.status").
			IsEqual(strconv.Itoa(status))
	}
}

func trackDataMatcher(status int, data any) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		r.Status(status).
			JSON().
			Schema(modelRef("trackData")).
			Object().
			ContainsValue(data)
	}
}

func tracksDataMatcher(status int, data any) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		r.Status(status).
			JSON().
			Schema(modelRef("tracksData")).
			Object().
			ContainsValue(data)
	}
}

func modelRef(modelName string) string {
	p, err := filepath.Abs("../../api/models.json")
	if err != nil {
		panic(err)
	}

	p, err = url.JoinPath("/", p)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("file://%s#/$def/%s", p, modelName)
}
