package api

import (
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/log"
)

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	encode(w, http.StatusMethodNotAllowed, errorResponse{
		Errors: []errorInfo{{
			Title:  "Method not allowed",
			Detail: fmt.Sprintf("The requested resource does not support a method '%s'", r.Method),
			Status: http.StatusMethodNotAllowed,
		}},
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	encode(w, http.StatusNotFound, errorResponse{
		Errors: []errorInfo{{
			Title:  "Resource not found",
			Detail: fmt.Sprintf("The requested path '%s' does not refer to a valid resource", r.URL.Path),
			Status: http.StatusNotFound,
		}},
	})
}

func internalError(msg string, err error, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Error(msg, err)

		encode(w, http.StatusInternalServerError, errorResponse{
			Errors: []errorInfo{{
				Title:  "Internal server error",
				Status: http.StatusInternalServerError,
			}},
		})
	}
}
