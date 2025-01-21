package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type TrackHandler struct {
	Store model.TrackStore
	Log   *log.Logger
}

func (h *TrackHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, h.getTrack)
}

func (h *TrackHandler) getTrack(r *http.Request) response {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return trackNotFound()
	}

	track, err := h.Store.TrackByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return trackNotFound()
		}

		h.serveError(r, err)
		return dataAccessError()
	}
	return dataFound(track)
}

func (h *TrackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, func(r *http.Request) response {
		tracks, err := h.Store.AllTracks()
		if err != nil {
			h.serveError(r, err)
			return dataAccessError()
		}
		return dataFound(tracks)
	})
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, h.createTrack)
}

func (h *TrackHandler) createTrack(r *http.Request) response {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req struct {
		Data model.TrackInfo `json:"data"`
	}

	if err := dec.Decode(&req); err != nil {
		h.serveError(r, err)
		return badRequest()
	}

	id, err := h.Store.CreateTrack(&req.Data)
	if err != nil {
		h.serveError(r, err)
		return resourceCreationError()
	}

	trackURI := r.URL.JoinPath(strconv.Itoa(id))
	track := model.Track{ID: id, Title: req.Data.Title}

	return dataCreated(trackURI.String(), &track)
}

func (h *TrackHandler) handleRequest(w http.ResponseWriter, r *http.Request, handler func(*http.Request) response) {
	if err := handler(r).write(w); err != nil {
		h.serveError(r, err)
	}
}

func (h *TrackHandler) serveError(r *http.Request, err error) {
	h.Log.WithContext(r.Context()).
		WithError(err).
		Error("Error serving the request")
}
