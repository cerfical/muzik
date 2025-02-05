package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/strutil"
)

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("The requested resource does not support a method '%s'", r.Method)
	encodeError(w, &apiError{
		Title:  "Method not allowed",
		Detail: details,
		Status: http.StatusMethodNotAllowed,
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("The requested path '%s' does not refer to a valid resource", r.URL.Path)
	encodeError(w, &apiError{
		Title:  "Resource not found",
		Detail: details,
		Status: http.StatusNotFound,
	})
}

func internalError(w http.ResponseWriter, _ *http.Request) {
	encodeError(w, &apiError{
		Title:  "Internal server error",
		Status: http.StatusInternalServerError,
	})
}

func badRequest(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if parseErr := (*parseError)(nil); errors.As(err, &parseErr) {
			encodeError(w, &apiError{
				Title:  "Request body is malformed",
				Detail: strutil.Capitalize(parseErr.Error()),
				Status: http.StatusBadRequest,
			})
		} else {
			internalError(w, r)
		}
	}
}
