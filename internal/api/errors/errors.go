package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/api"
	"github.com/cerfical/muzik/internal/strutil"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("the requested resource does not support a method '%s'", r.Method)
	e := Error{
		Title:  "method not allowed",
		Detail: details,
		Status: http.StatusMethodNotAllowed,
	}
	e.Write(w)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	details := fmt.Sprintf("the requested path '%s' does not refer to a valid resource", r.URL.Path)
	e := Error{
		Title:  "resource not found",
		Detail: details,
		Status: http.StatusNotFound,
	}
	e.Write(w)
}

func InternalError(w http.ResponseWriter, r *http.Request) {
	e := Error{
		Title:  "internal server error",
		Status: http.StatusInternalServerError,
	}
	e.Write(w)
}

func BadRequest(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if parseErr := (*api.ParseError)(nil); errors.As(err, &parseErr) {
			e := Error{
				Title:  "request body is malformed",
				Detail: parseErr.Error(),
				Status: http.StatusBadRequest,
			}
			e.Write(w)
		} else {
			InternalError(w, r)
		}
	}
}

type Error api.Error

func (e *Error) Write(w http.ResponseWriter) error {
	response := api.Response[struct{}]{
		Error:  (*api.Error)(e),
		Status: e.Status,
	}

	e.Title = strutil.Capitalize(e.Title)
	e.Detail = strutil.Capitalize(e.Detail)

	return response.Write(w)
}
