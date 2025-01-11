package tracks

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/model"
)

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

	enc := json.NewEncoder(wr)
	if err := enc.Encode(&resp); err != nil {
		log.Printf("writing response: %v", err)
	}
}
