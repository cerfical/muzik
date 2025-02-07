package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cerfical/muzik/internal/model"
)

const encodeMediaType = "application/json"

type trackDataResponse struct {
	Data *model.Track `json:"data"`
}

type tracksDataResponse struct {
	Data []model.Track `json:"data"`
}

type errorResponse struct {
	Errors []errorInfo `json:"errors"`
}

type errorInfo struct {
	Status int          `json:"status,string"`
	Title  string       `json:"title"`
	Detail string       `json:"detail,omitempty"`
	Source *errorSource `json:"source,omitempty"`
}

type errorSource struct {
	Header string `json:"header"`
}

type newTrackRequest struct {
	Data *model.Track `json:"data"`
}

func encode(w http.ResponseWriter, status int, r any) {
	w.Header().Set("Content-Type", encodeMediaType)
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(r)
}

func decode[T any](r io.Reader) (*T, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var req T
	if err := dec.Decode(&req); err != nil {
		if errMsg, ok := describeError(err); ok {
			return nil, &parseError{errMsg}
		}
		return nil, err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, &parseError{"The request body must contain a single JSON object"}
	}

	return &req, nil
}

func describeError(err error) (string, bool) {
	// Check for type errors
	if e := (&json.UnmarshalTypeError{}); errors.As(err, &e) {
		if e.Field != "" {
			return fmt.Sprintf("The request body contains an invalid value for the field '%s'", e.Field), true
		}
		return "The root element of the request body is invalid", true
	}

	// Check for generic syntax errors
	if e := (&json.SyntaxError{}); errors.As(err, &e) {
		return fmt.Sprintf("The request body has a syntax error at position %d", e.Offset), true
	}

	// Check for empty request body
	if errors.Is(err, io.EOF) {
		return "The request body must not be empty", true
	}

	// TODO: https://github.com/golang/go/issues/25956
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "The request body contains invalid JSON content", true
	}

	// TODO: https://github.com/golang/go/issues/29035
	errMsg := err.Error()
	if unknownFieldMsg := `json: unknown field "`; strings.HasPrefix(errMsg, unknownFieldMsg) {
		field := strings.TrimSuffix(strings.TrimPrefix(errMsg, unknownFieldMsg), `"`)
		return fmt.Sprintf("The request body contains an unknown field '%s'", field), true
	}

	// TODO: Provide more info about errors, maybe?
	if strings.HasPrefix(errMsg, "json: invalid use of ,string struct tag") {
		return "The request body contains fields with unexpected values", true
	}

	return "", false
}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}
