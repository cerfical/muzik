package handlers

import (
	"errors"
	"net/http"
	"strconv"

	apierrs "github.com/cerfical/muzik/internal/api/errors"

	"github.com/cerfical/muzik/internal/api"
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
		apierrs.NotFound(w, r)
		return
	}

	track, err := h.Store.TrackByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			apierrs.NotFound(w, r)
			return
		}

		h.Log.Error("failed to retrieve track data from storage", err)
		apierrs.InternalError(w, r)

		return
	}

	response := api.Response[*model.Track]{Data: &track, Status: http.StatusOK}
	response.Write(w)
}

func (h *Tracks) GetAll(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.Store.AllTracks()
	if err != nil {
		h.Log.Error("failed to retrieve tracks data from storage", err)
		apierrs.InternalError(w, r)
		return
	}

	response := api.Response[[]model.Track]{Data: &tracks, Status: http.StatusOK}
	response.Write(w)
}

func (h *Tracks) Create(w http.ResponseWriter, r *http.Request) {
	request, err := api.ReadRequest[model.Track](r.Body)
	if err != nil {
		apierrs.BadRequest(err)(w, r)
		return
	}

	if err := h.Store.CreateTrack(&request.Data); err != nil {
		h.Log.Error("failed to save track data to storage", err)
		apierrs.InternalError(w, r)
		return
	}

	location := r.URL.JoinPath(strconv.Itoa(request.Data.ID))
	w.Header().Set("Location", location.String())

	response := api.Response[model.Track]{Data: &request.Data, Status: http.StatusCreated}
	response.Write(w)
}
