package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/model"
)

type TrackHandler struct {
	Store model.TrackStore
}

func (h *TrackHandler) Get(wr http.ResponseWriter, req *http.Request) {
	trackID, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		writeErrorResponse(wr, http.StatusNotFound, "track not found")
		return
	}

	track, ok := h.Store.TrackByID(trackID)
	if !ok {
		writeErrorResponse(wr, http.StatusNotFound, "track not found")
		return
	}
	writeTrack(wr, track)
}

func (h *TrackHandler) GetAll(wr http.ResponseWriter, req *http.Request) {
	writeTracks(wr, h.Store.AllTracks())
}

func (h *TrackHandler) Create(wr http.ResponseWriter, req *http.Request) {
	info, err := readTrackInfo(req.Body)
	if err != nil {
		writeErrorResponse(wr, http.StatusBadRequest, err.Error())
		return
	}

	track := h.Store.CreateTrack(info)
	writeTrack(wr, track)
}

func readTrackInfo(r io.Reader) (*model.TrackInfo, error) {
	var req struct {
		Data *model.TrackInfo `json:"data"`
	}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	return req.Data, nil
}

func writeErrorResponse(wr http.ResponseWriter, status int, msg string) {
	r := struct {
		Errors []responseError `json:"errors"`
	}{[]responseError{{strconv.Itoa(status), msg}}}
	writeResponse(wr, status, r)
}

type responseError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

func writeTracks(wr http.ResponseWriter, data []*model.Track) {
	r := struct {
		Data []*model.Track `json:"data"`
	}{data}
	writeResponse(wr, http.StatusOK, r)
}

func writeTrack(wr http.ResponseWriter, data *model.Track) {
	r := struct {
		Data *model.Track `json:"data"`
	}{data}
	writeResponse(wr, http.StatusOK, r)
}

func writeResponse(wr http.ResponseWriter, status int, resp any) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(status)

	if err := json.NewEncoder(wr).Encode(&resp); err != nil {
		log.Printf("writing response: %v", err)
	}
}
