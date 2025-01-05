package serv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/stor"
)

type response struct {
	Status  int           `json:"-"`
	Message string        `json:"msg"`
	Data    []*stor.Track `json:"data"`
}

func (s *Server) index(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "static/index.html")
}

func (s *Server) trackByID(wr http.ResponseWriter, req *http.Request) {
	resp := s.lookupTrackByID(req.PathValue("id"))
	writeResponse(wr, resp)
}

func (s *Server) lookupTrackByID(trackID string) *response {
	trackNumID, err := strconv.Atoi(trackID)
	if err != nil {
		return malformedTrackID(trackID)
	}

	if trackNumID < 0 {
		return noTrackWithID(trackNumID)
	}

	track, ok := s.store.Get(trackNumID)
	if !ok {
		return noTrackWithID(trackNumID)
	}
	return tracksFound([]*stor.Track{track})
}

func (s *Server) listTracks(wr http.ResponseWriter, req *http.Request) {
	resp := tracksFound(s.store.GetAll())
	writeResponse(wr, resp)
}

func writeResponse(wr http.ResponseWriter, resp *response) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(resp.Status)

	// normalize the response so that the returned data is either empty ([]),
	// or populated with requested data items, but never null
	if resp.Data == nil {
		resp.Data = []*stor.Track{}
	}

	jsonenc := json.NewEncoder(wr)
	if err := jsonenc.Encode(&resp); err != nil {
		log.Printf("writing response: %v", err)
	}
}

func noTrackWithID(trackID int) *response {
	return makeErrorResponse(fmt.Sprintf("no track with such ID: %v", trackID))
}

func malformedTrackID(trackID string) *response {
	return makeErrorResponse(fmt.Sprintf("malformed track ID: %v", trackID))
}

func makeErrorResponse(msg string) *response {
	return &response{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

func tracksFound(data []*stor.Track) *response {
	return &response{
		Status: http.StatusOK,
		Data:   data,
	}
}
