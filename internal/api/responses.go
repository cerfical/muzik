package api

import (
	"net/http"

	"github.com/cerfical/muzik/internal/api/json"
	"github.com/cerfical/muzik/internal/log"
)

func badRequest(w http.ResponseWriter, r *http.Request, err string) {
	json.Error(w, r, err, http.StatusBadRequest)
}

func resourceNotExist(w http.ResponseWriter, r *http.Request) {
	json.Error(w, r, "resource doesn't exist", http.StatusNotFound)
}

func resourceReadError(w http.ResponseWriter, r *http.Request, err error) {
	internalError(w, r, "resource read failure", err)
}

func resourceCreationError(w http.ResponseWriter, r *http.Request, err error) {
	internalError(w, r, "resource creation failure", err)
}

func internalError(w http.ResponseWriter, r *http.Request, err string, e error) {
	log.FromRequest(r).
		WithError(e).Error("Internal error")

	json.Error(w, r, err, http.StatusInternalServerError)
}

func serveResource(w http.ResponseWriter, r *http.Request, data any) {
	json.Serve(w, r, data, http.StatusOK)
}

func resourceCreated(w http.ResponseWriter, r *http.Request, uri string, data any) {
	w.Header().Set("Location", uri)
	json.Serve(w, r, data, http.StatusCreated)
}
