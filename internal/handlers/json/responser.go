package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/strutil"
)

func NewResponser(w http.ResponseWriter, r *http.Request, log *log.Logger) *Responser {
	return &Responser{w, r, log}
}

// Responser provides a convenient interface for writing API responses.
type Responser struct {
	http.ResponseWriter
	request *http.Request
	log     *log.Logger
}

// ServeData writes resource data to the client.
func (r *Responser) ServeData(data any) {
	r.serveData(data, http.StatusOK)
}

// Created reports a successful resource creation.
func (r *Responser) Created(id int, data any) {
	resourceLocation := r.request.URL.JoinPath(strconv.Itoa(id))

	r.Header().Set("Location", resourceLocation.String())
	r.serveData(data, http.StatusCreated)
}

// NotFound reports a non-existent resource.
func (r *Responser) NotFound(id int) {
	r.error("track not found", fmt.Sprintf("no track with such ID: %d", id), http.StatusNotFound)
}

// BadRequest reports a malformed request body.
func (r *Responser) BadRequest(msg string) {
	r.error("bad request", msg, http.StatusBadRequest)
}

// RequestParseError reports an unexpected error in request parsing.
func (r *Responser) RequestParseError(err error) {
	r.internalError(
		"failed to parse the request",
		"something unexpected happened while parsing the request",
		err,
	)
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
	r.internalError("storage failure", detail, err)
}

func (r *Responser) internalError(title, detail string, err error) {
	r.logError(title, err)
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
	response.Error.Title = strutil.Capitalize(title)
	response.Error.Detail = strutil.Capitalize(detail)

	r.writeResponse(&response, code)
}

func (r *Responser) writeResponse(response any, status int) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)

	if err := json.NewEncoder(r).Encode(response); err != nil {
		r.logError("JSON encoding error", err)
	}
}

func (r *Responser) logError(msg string, err error) {
	if r.log != nil {
		r.log.Error(msg, err)
	}
}
