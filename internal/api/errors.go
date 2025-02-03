package api

import (
	"errors"
	"fmt"
	"net/http"
)

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("the requested resource does not support a method '%s'", r.Method)
	encodeError(w, &apiError{
		Title:  "method not allowed",
		Detail: details,
		Status: http.StatusMethodNotAllowed,
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("the requested path '%s' does not refer to a valid resource", r.URL.Path)
	encodeError(w, &apiError{
		Title:  "resource not found",
		Detail: details,
		Status: http.StatusNotFound,
	})
}

func internalError(w http.ResponseWriter, _ *http.Request) {
	encodeError(w, &apiError{
		Title:  "internal server error",
		Status: http.StatusInternalServerError,
	})
}

func badRequest(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if parseErr := (*parseError)(nil); errors.As(err, &parseErr) {
			encodeError(w, &apiError{
				Title:  "request body is malformed",
				Detail: parseErr.Error(),
				Status: http.StatusBadRequest,
			})
		} else {
			internalError(w, r)
		}
	}
}
