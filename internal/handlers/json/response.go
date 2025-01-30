package json

import (
	"encoding/json"
	"net/http"

	"github.com/cerfical/muzik/internal/strutil"
)

func Serve(w http.ResponseWriter, data any, code int) {
	r := struct {
		Data any `json:"data"`
	}{
		Data: data,
	}
	writeResponse(w, &r, code)
}

func Error(w http.ResponseWriter, title, detail string, code int) {
	var r struct {
		Error struct {
			Status int    `json:"status,string"`
			Title  string `json:"title"`
			Detail string `json:"detail,omitempty"`
		} `json:"error"`
	}

	r.Error.Status = code
	r.Error.Title = strutil.Capitalize(title)
	r.Error.Detail = strutil.Capitalize(detail)

	writeResponse(w, &r, code)
}

func writeResponse(w http.ResponseWriter, r any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// TODO: Ignoring errors?
	json.NewEncoder(w).Encode(r)
}
