package api

import (
	"encoding/json"
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

	track, ok := h.Store.TrackByID(id)
	if !ok {
		return trackNotFound()
	}
	return dataFound(track)
}

func (h *TrackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, func(r *http.Request) response {
		return dataFound(h.Store.AllTracks())
	})
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, h.createTrack)
}

func (h *TrackHandler) createTrack(r *http.Request) response {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req struct {
		Data *model.TrackInfo `json:"data"`
	}

	if err := dec.Decode(&req); err != nil {
		h.serveError(r, err)
		return badRequest()
	}

	track := h.Store.CreateTrack(req.Data)
	trackURI := r.URL.JoinPath(strconv.Itoa(track.ID))

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
