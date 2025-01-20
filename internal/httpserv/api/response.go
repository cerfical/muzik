package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func dataFound(data any) response {
	return &dataResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}

func dataCreated(uri string, data any) response {
	return &dataResponse{
		Status: http.StatusCreated,
		Header: []string{"Location", uri},
		Data:   data,
	}
}

func trackNotFound() response {
	return &errorResponse{
		Status: http.StatusNotFound,
		Error:  "track not found",
	}
}

func badRequest() response {
	return &errorResponse{
		Status: http.StatusBadRequest,
		Error:  "failed to parse the request",
	}
}

type response interface {
	write(w http.ResponseWriter) error
}

type errorResponse struct {
	Status int      `json:"-"`
	Header []string `json:"-"`
	Error  string   `json:"-"`
}

func (r *errorResponse) write(w http.ResponseWriter) error {
	resp := struct {
		Errors []responseError `json:"errors"`
	}{[]responseError{{
		Status: strconv.Itoa(r.Status),
		Title:  r.Error,
	}}}

	return writeResponse(resp, r.Status, r.Header, w)
}

type responseError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

type dataResponse struct {
	Status int      `json:"-"`
	Header []string `json:"-"`
	Data   any      `json:"data"`
}

func (r *dataResponse) write(w http.ResponseWriter) error {
	return writeResponse(r, r.Status, r.Header, w)
}

func writeResponse(resp any, status int, header []string, w http.ResponseWriter) error {
	for i := 0; i < len(header); i += 2 {
		w.Header().Set(header[i], header[i+1])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(resp)
}
