package api

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Data   *T     `json:"data,omitempty"`
	Error  *Error `json:"error,omitempty"`
	Status int    `json:"-"`
}

type Error struct {
	Status int    `json:"status,string"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (r *Response[T]) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)

	return json.NewEncoder(w).Encode(r)
}
