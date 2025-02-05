package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const encodeMediaType = "application/json"

type apiError struct {
	Status int          `json:"status,string"`
	Title  string       `json:"title"`
	Detail string       `json:"detail,omitempty"`
	Source *errorSource `json:"source,omitempty"`
}

type errorSource struct {
	Header string `json:"header"`
}

func decodeData[T any](r io.Reader) (T, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var empty T
	var request struct {
		Data struct {
			Attrs T `json:"attributes"`
		} `json:"data"`
	}

	if err := dec.Decode(&request); err != nil {
		if errMsg, ok := describeError(err); ok {
			return empty, &parseError{errMsg}
		}
		return empty, err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return empty, &parseError{"request body must contain a single JSON object"}
	}

	return request.Data.Attrs, nil
}

func describeError(err error) (string, bool) {
	// Check for type errors
	if e := (&json.UnmarshalTypeError{}); errors.As(err, &e) {
		if e.Field != "" {
			return fmt.Sprintf("invalid value for a field '%s'", e.Field), true
		}
		return "root element is invalid", true
	}

	// Check for generic syntax errors
	if e := (&json.SyntaxError{}); errors.As(err, &e) {
		return fmt.Sprintf("syntax error at position %d", e.Offset), true
	}

	// Check for empty request body
	if errors.Is(err, io.EOF) {
		return "request body must not be empty", true
	}

	// TODO: https://github.com/golang/go/issues/25956
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "request body is not valid JSON", true
	}

	// TODO: https://github.com/golang/go/issues/29035
	errMsg := err.Error()
	if unknownFieldMsg := `json: unknown field "`; strings.HasPrefix(errMsg, unknownFieldMsg) {
		field := strings.TrimSuffix(strings.TrimPrefix(errMsg, unknownFieldMsg), `"`)
		return fmt.Sprintf("unknown field '%s'", field), true
	}

	// TODO: Provide more info about errors, maybe?
	if strings.HasPrefix(errMsg, "json: invalid use of ,string struct tag") {
		return "some fields are expected to be strings, but are not", true
	}

	return "", false
}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}

func encodeData[T any](w http.ResponseWriter, data T, status int) error {
	response := struct {
		Data T `json:"data"`
	}{
		Data: data,
	}

	return encodeJSON(w, &response, status)
}

func encodeError(w http.ResponseWriter, e *apiError) error {
	response := struct {
		Errors []*apiError `json:"errors"`
	}{
		Errors: []*apiError{e},
	}

	return encodeJSON(w, &response, e.Status)
}

func encodeJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", encodeMediaType)
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
