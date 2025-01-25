package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/handlers/json"
	"github.com/cerfical/muzik/internal/model"
)

type Tracks struct {
	Store model.TrackStore
}

func (h *Tracks) Get(w http.ResponseWriter, r *http.Request) {
	responser := json.NewResponser(w, r)

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
	responser := json.NewResponser(w, r)

	tracks, err := h.Store.AllTracks()
	if err != nil {
		responser.StorageReadError(err)
		return
	}

	responser.ServeData(tracks)
}

func (h *Tracks) Create(w http.ResponseWriter, r *http.Request) {
	responser := json.NewResponser(w, r)

	var attrs model.TrackInfo
	if err := json.ParseData(r.Body, &attrs); err != nil {
		if parseErr := (*json.ParseError)(nil); errors.As(err, &parseErr) {
			responser.MalformedRequest(parseErr.Error())
		} else {
			responser.RequestParseError(err)
		}
		return
	}

	id, err := h.Store.CreateTrack(&attrs)
	if err != nil {
		responser.StorageWriteError(err)
		return
	}

	track := model.Track{ID: id, Title: attrs.Title}
	responser.Created(id, &track)
}
