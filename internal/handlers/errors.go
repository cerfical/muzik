package handlers

import (
	"net/http"

	"github.com/cerfical/muzik/internal/handlers/json"
)

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
