package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func WriteError(wr http.ResponseWriter, status int, msg string) {
	resp := errorResponse{
		Errors: []responseError{
			{strconv.Itoa(status), msg},
		},
	}
	writeResponse(wr, status, resp)
}

func WriteDataItems[T any](wr http.ResponseWriter, data []T) {
	r := dataItemsResponse[T]{Data: data}
	writeResponse(wr, http.StatusOK, r)
}

func WriteDataItem[T any](wr http.ResponseWriter, data T) {
	r := dataItemResponse[T]{Data: data}
	writeResponse(wr, http.StatusOK, r)
}

func writeResponse(wr http.ResponseWriter, status int, resp any) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(status)

	e := json.NewEncoder(wr)
	if err := e.Encode(&resp); err != nil {
		log.Printf("writing response: %v", err)
	}
}
