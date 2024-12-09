package httpserv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/storage"
)

type response struct {
	Message string           `json:"msg"`
	Data    []*storage.Track `json:"data"`
}

func (s *Server) getTrack(wr http.ResponseWriter, req *http.Request) {
	resp, status := s.lookupTrackByID(req.PathValue("id"))

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(status)

	jsonenc := json.NewEncoder(wr)
	if err := jsonenc.Encode(&resp); err != nil {
		log.Printf("writing response: %v", err)
	}
}

func (s *Server) lookupTrackByID(trackID string) (response, int) {
	trackNumID, err := strconv.Atoi(trackID)
	if err != nil {
		return malformedTrackID(trackID), http.StatusNotFound
	}

	if trackNumID < 0 {
		return noTrackWithID(trackNumID), http.StatusNotFound
	}

	track, ok := s.store.Get(trackNumID)
	if !ok {
		return noTrackWithID(trackNumID), http.StatusNotFound
	}

	return response{
		Message: "",
		Data:    []*storage.Track{track},
	}, http.StatusOK
}

func noTrackWithID(trackID int) response {
	return makeErrorResponse(fmt.Sprintf("no track with such ID: %v", trackID))
}

func malformedTrackID(trackID string) response {
	return makeErrorResponse(fmt.Sprintf("malformed track ID: %v", trackID))
}

func makeErrorResponse(msg string) response {
	return response{
		Message: msg,
		Data:    []*storage.Track{},
	}
}
