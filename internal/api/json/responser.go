package json

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/log"
)

func NewResponser(w http.ResponseWriter, r *http.Request) *Responser {
	return &Responser{w, r}
}

// Responser provides a convenient interface for writing API responses.
type Responser struct {
	rw  http.ResponseWriter
	req *http.Request
}

// ServeData writes resource data to the client.
func (r *Responser) ServeData(data any) {
	r.serveData(data, http.StatusOK)
}

// Created reports a successful resource creation.
func (r *Responser) Created(id int, data any) {
	resourceLocation := r.req.URL.JoinPath(strconv.Itoa(id))

	r.rw.Header().Set("Location", resourceLocation.String())
	r.serveData(data, http.StatusCreated)
}

// NotFound reports a non-existent resource.
func (r *Responser) NotFound() {
	r.error("track not found", "track with such ID does not exist", http.StatusNotFound)
}

// MalformedRequest reports a malformed request body.
func (r *Responser) MalformedRequest(msg string) {
	r.error("malformed request body", msg, http.StatusBadRequest)
}

// RequestParseError reports an unexpected error in request parsing.
func (r *Responser) RequestParseError(err error) {
	r.logError("Failed to parse the request body", err)
	r.internalError("request parsing failure", "something unexpected happened while parsing the request body")
}

// StorageWriteError reports a write error to the storage.
func (r *Responser) StorageWriteError(err error) {
	r.storageError("unable to store track data", err)
}

// StorageReadError reports a read error from the storage.
func (r *Responser) StorageReadError(err error) {
	r.storageError("unable to get track data", err)
}

func (r *Responser) storageError(detail string, err error) {
	r.logError("Storage failure", err)
	r.internalError("storage failure", detail)
}

func (r *Responser) internalError(title, detail string) {
	r.error(title, detail, http.StatusInternalServerError)
}

func (r *Responser) serveData(data any, code int) {
	response := struct {
		Data any `json:"data"`
	}{
		Data: data,
	}

	r.writeResponse(&response, code)
}

func (r *Responser) error(title, detail string, code int) {
	var response struct {
		Error struct {
			Status int    `json:"status,string"`
			Title  string `json:"title,omitempty"`
			Detail string `json:"detail,omitempty"`
		} `json:"error"`
	}

	response.Error.Status = code
	response.Error.Title = title
	response.Error.Detail = detail

	r.writeResponse(&response, code)
}

func (r *Responser) writeResponse(response any, status int) {
	r.rw.Header().Set("Content-Type", "application/json")
	r.rw.WriteHeader(status)

	if err := json.NewEncoder(r.rw).Encode(response); err != nil {
		r.logError("JSON encoding error", err)
	}
}

func (r *Responser) logError(msg string, err error) {
	log.FromRequest(r.req).
		WithError(err).Error(msg)
}
