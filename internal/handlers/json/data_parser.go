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
		if errMsg, ok := describeError(err); ok {
			return &ParseError{errMsg}
		}
		return err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return &ParseError{"request body must contain a single JSON object"}
	}

	return nil
}

func describeError(err error) (string, bool) {
	// Check for type errors
	if e := (&json.UnmarshalTypeError{}); errors.As(err, &e) {
		if e.Field != "" {
			return fmt.Sprintf("invalid value for a field: %s", e.Field), true
		}
		return "root element is invalid", true
	}

	// Check for generic syntax errors
	if e := (&json.SyntaxError{}); errors.As(err, &e) {
		return fmt.Sprintf("syntax error at offset: %d", e.Offset), true
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
		return fmt.Sprintf("unknown field: %s", field), true
	}

	// TODO: Provide more info about errors, maybe?
	if strings.HasPrefix(errMsg, "json: invalid use of ,string struct tag") {
		return "some fields are expected to be strings, but are not", true
	}

	return "", false
}

// ParseError describes syntax errors that occur during [ParseData].
type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}
