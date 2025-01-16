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
	trackID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, http.StatusNotFound, "track not found")
		return
	}

	track, ok := h.Store.TrackByID(trackID)
	if !ok {
		h.writeError(w, r, http.StatusNotFound, "track not found")
		return
	}
	h.writeTrack(w, r, http.StatusOK, track)
}

func (h *TrackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.writeTracks(w, r, http.StatusOK, h.Store.AllTracks())
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var trackInfo struct {
		Data *model.TrackInfo `json:"data"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&trackInfo); err != nil {
		h.writeError(w, r, http.StatusBadRequest, "failed to decode request body")
		return
	}

	track := h.Store.CreateTrack(trackInfo.Data)

	trackURI := r.URL.JoinPath(strconv.Itoa(track.ID))
	w.Header().Set("Location", trackURI.String())
	h.writeTrack(w, r, http.StatusCreated, track)
}

func (h *TrackHandler) writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	resp := struct {
		Errors []responseError `json:"errors"`
	}{[]responseError{{strconv.Itoa(status), msg}}}

	h.writeResponse(w, r, status, resp)
}

type responseError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

func (h *TrackHandler) writeTracks(w http.ResponseWriter, r *http.Request, status int, data []*model.Track) {
	resp := struct {
		Data []*model.Track `json:"data"`
	}{data}
	h.writeResponse(w, r, status, resp)
}

func (h *TrackHandler) writeTrack(w http.ResponseWriter, r *http.Request, status int, data *model.Track) {
	resp := struct {
		Data *model.Track `json:"data"`
	}{data}
	h.writeResponse(w, r, status, resp)
}

func (h *TrackHandler) writeResponse(w http.ResponseWriter, r *http.Request, status int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.Log.WithContext(r.Context()).
			WithError(err).
			Error("Failed to encode response body")
	}
}
