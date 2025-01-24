package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/api/json"
	"github.com/cerfical/muzik/internal/model"
)

type TrackHandler struct {
	Store model.TrackStore
}

func (h *TrackHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		resourceNotExist(w, r)
		return
	}

	track, err := h.Store.TrackByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			resourceNotExist(w, r)
			return
		}

		resourceReadError(w, r, err)
		return
	}

	serveResource(w, r, track)
}

func (h *TrackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.Store.AllTracks()
	if err != nil {
		resourceReadError(w, r, err)
		return
	}

	serveResource(w, r, tracks)
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var attrs model.TrackInfo
	if err := json.Parse(r.Body, &attrs); err != nil {
		badRequest(w, r, "failed to parse the request body")
		return
	}

	id, err := h.Store.CreateTrack(&attrs)
	if err != nil {
		resourceCreationError(w, r, err)
		return
	}

	trackURI := r.URL.JoinPath(strconv.Itoa(id))
	track := model.Track{ID: id, Title: attrs.Title}

	resourceCreated(w, r, trackURI.String(), &track)
}
