package handlers_test

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/cerfical/muzik/internal/api/handlers"
	"github.com/cerfical/muzik/internal/mocks"
	"github.com/cerfical/muzik/internal/model"
	"github.com/stretchr/testify/suite"

	"github.com/gavv/httpexpect/v2"
)

func TestTracks(t *testing.T) {
	suite.Run(t, new(TracksTest))
}

type TracksTest struct {
	suite.Suite

	store  *mocks.TrackStore
	expect *httpexpect.Expect
}

func (t *TracksTest) SetupSubTest() {
	t.store = mocks.NewTrackStore(t.T())
	tracks := handlers.Tracks{Store: t.store}

	router := http.NewServeMux()
	router.HandleFunc("GET /{id}", tracks.Get)
	router.HandleFunc("GET /", tracks.GetAll)
	router.HandleFunc("POST /", tracks.Create)

	client := http.Client{Transport: httpexpect.NewBinder(router)}
	t.expect = httpexpect.WithConfig(httpexpect.Config{
		Client:   &client,
		TestName: t.T().Name(),
		Reporter: httpexpect.NewAssertReporter(t.T()),
	})
}

func (t *TracksTest) TestTracks_Get() {
	track := model.Track{ID: 7, Title: "Example Track"}
	tests := []struct {
		name      string
		id        any
		storeData *model.Track
		storeErr  error
		matcher   func(*httpexpect.Response)
	}{
		{"ok", track.ID, &track, nil, dataMatcher(http.StatusOK, &track)},
		{"not_found", track.ID, nil, model.ErrNotFound, errorMatcher(http.StatusNotFound)},
		{"invalid_id", "invalid_id", nil, nil, errorMatcher(http.StatusNotFound)},
		{"storage_fail", track.ID, nil, errors.New(""), errorMatcher(http.StatusInternalServerError)},
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

func (t *TracksTest) TestTracks_GetAll() {
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
		{"ok", tracks, nil, dataMatcher(http.StatusOK, tracks)},
		{"empty", empty, nil, dataMatcher(http.StatusOK, empty)},
		{"storage_fail", nil, errors.New(""), errorMatcher(http.StatusInternalServerError)},
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

func (t *TracksTest) TestTracks_Create() {
	track := model.Track{ID: 7, Title: "Example Track"}
	body := struct {
		Data *model.Track `json:"data"`
	}{
		Data: &track,
	}

	tests := []struct {
		name      string
		storeData *model.Track
		storeErr  error
		body      any
		location  string
		matcher   func(*httpexpect.Response)
	}{
		{"ok", &track, nil, body, "/7", dataMatcher(http.StatusCreated, &track)},
		{"bad_request", nil, nil, "{}", "", errorMatcher(http.StatusBadRequest)},
		{"storage_fail", &track, errors.New(""), body, "", errorMatcher(http.StatusInternalServerError)},
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
		r.Status(status).JSON().
			Object().
			Value("error").
			Object().
			ContainsKey("title").
			ContainsKey("detail").
			HasValue("status", strconv.Itoa(status))
	}
}

func dataMatcher(status int, data any) func(*httpexpect.Response) {
	return func(r *httpexpect.Response) {
		r.Status(status).JSON().
			Object().
			HasValue("data", data)
	}
}
