package handlers

import (
	"errors"
	"fmt"
	"mime"
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
	if !validateContentType(w, r) {
		return
	}

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

func validateContentType(w http.ResponseWriter, r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	if isJSONContent(contentType) {
		return true
	}

	if h, ok := acceptHeaderForMethod(r.Method); ok {
		w.Header().Set(h, jsonMediaType)
	}

	json.Error(w,
		"content body has an unknown media type",
		fmt.Sprintf("unexpected content type '%s', only '%s' is supported", contentType, jsonMediaType),
		http.StatusUnsupportedMediaType,
	)

	return false
}

func isJSONContent(contentType string) bool {
	if len(contentType) == 0 {
		// Ignore empty Content-Type headers
		return true
	}

	contentType, params, err := mime.ParseMediaType(contentType)
	if err != nil || contentType != jsonMediaType {
		return false
	}

	if n := len(params); n > 0 {
		// The only allowed media parameter is charset=utf-8, redundant for application/json
		if n == 1 && params["charset"] == "utf-8" {
			return true
		}
		return false
	}

	return true
}

func acceptHeaderForMethod(method string) (string, bool) {
	switch method {
	case http.MethodPatch:
		return "Accept-Patch", true
	case http.MethodPost:
		return "Accept-Post", true
	default:
		return "", false
	}
}

const (
	jsonMediaType = "application/json"
)
