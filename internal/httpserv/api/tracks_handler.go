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

	track, err := h.store.GetTrack(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			notFound(w, r)
			return
		}
		internalError("Failed to read track data from persistent storage", err, h.log)(w, r)
		return
	}

	encode(w, http.StatusOK, trackDataResponse{
		Data: track,
	})
}

func (h *tracksHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.store.GetTracks(r.Context())
	if err != nil {
		internalError("Failed to read tracks data from persistent storage", err, h.log)(w, r)
		return
	}

	encode(w, http.StatusOK, tracksDataResponse{
		Data: tracks,
	})
}

func (h *tracksHandler) create(w http.ResponseWriter, r *http.Request) {
	newTrack, err := decode[newTrackRequest](r.Body)
	if err != nil {
		if parseErr := (*parseError)(nil); errors.As(err, &parseErr) {
			encode(w, http.StatusBadRequest, errorResponse{
				Errors: []errorInfo{{
					Title:  "The request body is malformed",
					Detail: parseErr.Error(),
					Status: http.StatusBadRequest,
				}},
			})
		} else {
			internalError("Parsing of the request body was interrupted due to an unexpected error", err, h.log)(w, r)
		}
		return
	}

	track, err := h.store.CreateTrack(r.Context(), (*model.TrackAttrs)(&newTrack.Data.Attrs))
	if err != nil {
		internalError("Failed to save track data to persistent storage", err, h.log)(w, r)
		return
	}

	location := r.URL.JoinPath(strconv.Itoa(track.ID))
	w.Header().Set("Location", location.String())

	encode(w, http.StatusCreated, trackDataResponse{
		Data: track,
	})
}

func (h *tracksHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		notFound(w, r)
		return
	}

	if err := h.store.DeleteTrack(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			notFound(w, r)
		} else {
			internalError("Failed to delete track data from persistent storage", err, h.log)(w, r)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
