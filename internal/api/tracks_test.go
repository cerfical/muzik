package api_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cerfical/muzik/internal/model"
)

func (t *APITest) TestTracks_Get() {
	track := model.Track{ID: 7, Title: "Example Track"}
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
					TrackByID(id).
					Return(&track, test.storeErr)
			}

			e := t.expect.GET("/{id}", test.id)
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("trackData", test.status, &track))
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
		{"ok_200",
			[]model.Track{
				{ID: 7, Title: "Example Track #1"},
				{ID: 2, Title: "Example Track #2"},
			},
			nil,
			http.StatusOK,
		},
		{"ok_200_empty", []model.Track{}, nil, http.StatusOK},

		{"500_storage_fail", nil, errors.New(""), http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			t.store.EXPECT().
				AllTracks().
				Return(test.storeData, test.storeErr)

			e := t.expect.GET("/")
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("tracksData", test.status, test.storeData))
			} else {
				e = e.WithMatcher(isErrorResponse(test.status))
			}
			e.Expect()
		})
	}
}

func (t *APITest) TestTracks_Create() {
	trackData := struct {
		Data model.Track `json:"data"`
	}{
		Data: model.Track{ID: 7, Title: "Example Track"},
	}

	tests := []struct {
		name string

		storeData *model.Track
		storeErr  error
		body      any
		location  string
		status    int
	}{
		{"ok_201",
			&trackData.Data,
			nil,
			&trackData,
			"/api/tracks/7",
			http.StatusCreated,
		},

		{"fail_400_bad_request", nil, nil, "{}", "", http.StatusBadRequest},
		{"fail_500_storage_fail", &trackData.Data, errors.New(""), &trackData, "", http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.storeData != nil {
				t.store.EXPECT().
					CreateTrack(test.storeData).
					Return(test.storeErr)
			}

			e := t.expect.POST("/")
			if strings.HasPrefix(test.name, "ok") {
				e = e.WithMatcher(isDataResponse("trackData", test.status, test.storeData))
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
