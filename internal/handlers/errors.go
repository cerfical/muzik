package handlers

import (
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/handlers/json"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"method not allowed",
		fmt.Sprintf("the requested resource does not support a method '%s'", r.Method),
		http.StatusMethodNotAllowed,
	)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"resource not found",
		fmt.Sprintf("the requested path '%s' does not refer to a valid resource", r.URL.Path),
		http.StatusNotFound,
	)
}

func InternalError(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"internal server error",
		"something very bad happened on the server side",
		http.StatusInternalServerError,
	)
}
