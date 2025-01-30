package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/handlers/json"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type Tracks struct {
	Store model.TrackStore
	Log   *log.Logger
}

func (h *Tracks) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		NotFound(w, r)
		return
	}

	track, err := h.Store.TrackByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			NotFound(w, r)
			return
		}

		h.Log.Error("error reading from the database", err)
		InternalError(w, r)

		return
	}

	json.Serve(w, track, http.StatusOK)
}

func (h *Tracks) GetAll(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.Store.AllTracks()
	if err != nil {
		h.Log.Error("error reading from the database", err)
		InternalError(w, r)
		return
	}

	json.Serve(w, tracks, http.StatusOK)
}

func (h *Tracks) Create(w http.ResponseWriter, r *http.Request) {
	var track model.Track
	if err := json.Read(r.Body, &track); err != nil {
		if parseErr := (*json.ParseError)(nil); errors.As(err, &parseErr) {
			json.Error(w, "request body is malformed", parseErr.Error(), http.StatusBadRequest)
		} else {
			h.Log.Error("failed to read the request", err)
			InternalError(w, r)
		}
		return
	}

	if err := h.Store.CreateTrack(&track); err != nil {
		h.Log.Error("error writing to the database", err)
		InternalError(w, r)
		return
	}

	location := r.URL.JoinPath(strconv.Itoa(track.ID))
	w.Header().Set("Location", location.String())

	json.Serve(w, &track, http.StatusCreated)
}
