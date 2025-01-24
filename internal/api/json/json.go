package json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cerfical/muzik/internal/log"
)

// Error writes a JSON error response to the client.
func Error(w http.ResponseWriter, r *http.Request, err string, code int) {
	var response struct {
		Error struct {
			Status int    `json:"status,string"`
			Title  string `json:"title"`
		} `json:"error"`
	}

	response.Error.Status = code
	response.Error.Title = err

	writeResponse(w, r, &response, code)
}

// Serve serves a JSON resource to the client.
func Serve(w http.ResponseWriter, r *http.Request, data any, code int) {
	response := dataResponse{Data: data}
	writeResponse(w, r, &response, code)
}

func writeResponse(w http.ResponseWriter, r *http.Request, response any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.FromRequest(r).
			WithError(err).Error("JSON encoding error")
	}
}

// Parse parses a JSON resource into a struct.
func Parse(r io.Reader, data any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	request := dataResponse{Data: data}
	if err := dec.Decode(&request); err != nil {
		return err
	}
	return nil
}

type dataResponse struct {
	Data any `json:"data"`
}
