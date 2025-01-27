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
	responser := json.NewResponser(w, r, h.Log)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		responser.NotFound()
		return
	}

	track, err := h.Store.TrackByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			responser.NotFound()
			return
		}

		responser.StorageReadError(err)
		return
	}

	responser.ServeData(track)
}

func (h *Tracks) GetAll(w http.ResponseWriter, r *http.Request) {
	responser := json.NewResponser(w, r, h.Log)

	tracks, err := h.Store.AllTracks()
	if err != nil {
		responser.StorageReadError(err)
		return
	}

	responser.ServeData(tracks)
}

func (h *Tracks) Create(w http.ResponseWriter, r *http.Request) {
	responser := json.NewResponser(w, r, h.Log)

	var track model.Track
	if err := json.ParseData(r.Body, &track); err != nil {
		if parseErr := (*json.ParseError)(nil); errors.As(err, &parseErr) {
			responser.MalformedRequest(parseErr.Error())
		} else {
			responser.RequestParseError(err)
		}
		return
	}

	if err := h.Store.CreateTrack(&track); err != nil {
		responser.StorageWriteError(err)
		return
	}

	responser.Created(track.ID, &track)
}
