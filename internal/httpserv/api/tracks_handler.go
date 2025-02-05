package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type tracksHandler struct {
	store model.TrackStore
	log   *log.Logger
}

func (h *tracksHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		notFound(w, r)
		return
	}

	track, err := h.store.TrackByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			notFound(w, r)
			return
		}

		h.log.Error("failed to retrieve track data from storage", err)
		internalError(w, r)

		return
	}

	encodeData(w, &track, http.StatusOK)
}

func (h *tracksHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.store.AllTracks(r.Context())
	if err != nil {
		h.log.Error("failed to retrieve tracks data from storage", err)
		internalError(w, r)
		return
	}

	encodeData(w, tracks, http.StatusOK)
}

func (h *tracksHandler) create(w http.ResponseWriter, r *http.Request) {
	track, err := decodeData[*model.Track](r.Body)
	if err != nil {
		badRequest(err)(w, r)
		return
	}

	if err := h.store.CreateTrack(r.Context(), track); err != nil {
		h.log.Error("failed to save track data to storage", err)
		internalError(w, r)
		return
	}

	location := r.URL.JoinPath(strconv.Itoa(track.ID))
	w.Header().Set("Location", location.String())

	encodeData(w, track, http.StatusCreated)
}
