package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ParseData parses a JSON resource into a struct.
func ParseData(r io.Reader, data any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	request := struct {
		Data any `json:"data"`
	}{
		Data: data,
	}

	if err := dec.Decode(&request); err != nil {
		if errMsg := describeError(err); errMsg != "" {
			return &ParseError{errMsg}
		}
		return err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return &ParseError{"request body must contain a single JSON object"}
	}

	return nil
}

func describeError(err error) string {
	// Check for type errors
	if e := (&json.UnmarshalTypeError{}); errors.As(err, &e) {
		if e.Field != "" {
			return fmt.Sprintf("request body contains an invalid value for the field %s", e.Field)
		}
		return "request body contains an invalid value for the root element"
	}

	// Check for generic syntax errors
	if e := (&json.SyntaxError{}); errors.As(err, &e) {
		return fmt.Sprintf("request body has a syntax error at position %d", e.Offset)
	}

	// Check for empty request body
	if errors.Is(err, io.EOF) {
		return "request body must not be empty"
	}

	// TODO: https://github.com/golang/go/issues/25956
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "request body contains malformed JSON"
	}

	// TODO: https://github.com/golang/go/issues/29035
	if field := strings.TrimPrefix(err.Error(), "json: unknown field \""); field != err.Error() {
		field = strings.TrimSuffix(field, "\"")
		return fmt.Sprintf("request body contains unknown field %s", field)
	}

	return ""
}

// ParseError describes syntax errors that occur during [ParseData].
type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}
