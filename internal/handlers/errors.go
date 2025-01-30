package handlers

import (
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/handlers/json"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"request method is not allowed",
		fmt.Sprintf("the endpoint does not support a method '%s'", r.Method),
		http.StatusMethodNotAllowed,
	)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"the requested resource does not exist",
		"",
		http.StatusNotFound,
	)
}

func InternalError(w http.ResponseWriter, r *http.Request) {
	json.Error(w,
		"internal server failure",
		"",
		http.StatusInternalServerError,
	)
}
