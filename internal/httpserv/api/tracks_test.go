package api_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cerfical/muzik/internal/model"
	"github.com/stretchr/testify/mock"
)

var sampleTrack = model.Track{
	ID: 7,
	Attrs: model.TrackAttrs{
		Title: "Example Track",
	},
}

func (t *APITest) TestTracks_Get() {
	tests := []struct {
		name string

		id       string
		storeErr error
		status   int
	}{
		{"ok_200", "7", nil, http.StatusOK},

		{"fail_404_not_found", "7", model.ErrNotFound, http.StatusNotFound},
		{"fail_404_invalid_id", "abcd-efgh", nil, http.StatusNotFound},
		{"fail_500_storage_fail", "7", errors.New(""), http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if id, err := strconv.Atoi(test.id); err == nil {
				t.store.EXPECT().
					TrackByID(mock.Anything, id).
					Return(&sampleTrack, test.storeErr)
			}

			e := t.expect.GET("/{id}", test.id)
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("TrackResource", test.status, &sampleTrack))
			} else {
				e = e.WithMatcher(isErrorResponse(test.status))
			}
			e.Expect()
		})
	}
}

func (t *APITest) TestTracks_GetAll() {
	tests := []struct {
		name string

		storeData []model.Track
		storeErr  error
		status    int
	}{
		{"ok_200", []model.Track{sampleTrack, sampleTrack}, nil, http.StatusOK},
		{"ok_200_empty", []model.Track{}, nil, http.StatusOK},

		{"500_storage_fail", nil, errors.New(""), http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			t.store.EXPECT().
				AllTracks(mock.Anything).
				Return(test.storeData, test.storeErr)

			e := t.expect.GET("/")
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("TracksResource", test.status, test.storeData))
			} else {
				e = e.WithMatcher(isErrorResponse(test.status))
			}
			e.Expect()
		})
	}
}

func (t *APITest) TestTracks_Create() {
	var trackData struct {
		Data struct {
			Attrs model.TrackAttrs `json:"attributes"`
		} `json:"data"`
	}
	trackData.Data.Attrs = sampleTrack.Attrs

	tests := []struct {
		name string

		storeData *model.Track
		storeErr  error
		body      any
		location  string
		status    int
	}{
		{"ok_201", &sampleTrack, nil, &trackData, "/api/tracks/7", http.StatusCreated},

		{"fail_400_bad_request", nil, nil, "{}", "", http.StatusBadRequest},
		{"fail_500_storage_fail", &sampleTrack, errors.New(""), &trackData, "", http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.storeData != nil || test.storeErr != nil {
				t.store.EXPECT().
					CreateTrack(mock.Anything, mock.Anything).
					Return(test.storeData, test.storeErr)
			}

			e := t.expect.POST("/")
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("TrackResource", test.status, test.storeData))
			} else {
				e = e.WithMatcher(isErrorResponse(test.status))
			}

			e.WithJSON(test.body).
				Expect().
				Header("Location").
				IsEqual(test.location)
		})
	}
}
