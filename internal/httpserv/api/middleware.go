package api

import (
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/cerfical/muzik/internal/log"
	"github.com/pkg/errors"
)

// accepts checks Accept header for the presence of the specified media type.
func accepts(mediaType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if checkAcceptHeader(mediaType, r.Header.Get("Accept")) {
				next.ServeHTTP(w, r)
				return
			}

			encode(w, http.StatusNotAcceptable, errorResponse{
				Errors: []errorInfo{{
					Title:  "Media type is not acceptable",
					Detail: fmt.Sprintf("The only acceptable media type is '%s'", mediaType),
					Status: http.StatusNotAcceptable,
					Source: &errorSource{
						Header: "Accept",
					},
				}},
			})
		}
	}
}

func checkAcceptHeader(supportedType, acceptHeader string) bool {
	if acceptHeader == "" {
		// Ignore empty Accept headers
		return true
	}

	supMain, supSub := splitMediaType(supportedType)
	accTypes := strings.Split(acceptHeader, ",")

	for _, accType := range accTypes {
		accType, params, err := mime.ParseMediaType(accType)
		if err != nil {
			continue
		}

		if len(params) > 0 {
			// Q-values is the only allowed media type parameter
			val, ok := params["q"]
			if len(params) != 1 || !ok {
				continue
			}

			// Q-value is invalid, or explicitly set to 0
			if valNum, err := strconv.ParseFloat(val, 64); err != nil || valNum == 0 {
				continue
			}
		}

		accMain, accSub := splitMediaType(accType)
		if (accMain != supMain && accMain != "*") || (accSub != supSub && accSub != "*") {
			continue
		}

		return true
	}

	return false
}

func splitMediaType(mediaType string) (string, string) {
	mainType, subType, _ := strings.Cut(mediaType, "/")
	return mainType, subType
}

// hasContentType checks Content-Type for the presence of the specified media type.
func hasContentType(mediaType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if !hasContentBody(r) || checkContentType(contentType, mediaType) {
				next.ServeHTTP(w, r)
				return
			}

			if h, ok := acceptHeaderForMethod(r.Method); ok {
				w.Header().Set(h, mediaType)
			}

			encode(w, http.StatusUnsupportedMediaType, errorResponse{
				Errors: []errorInfo{{
					Title:  "Media type is unsupported",
					Detail: fmt.Sprintf("Unexpected content type '%s', only '%s' is allowed", contentType, mediaType),
					Status: http.StatusUnsupportedMediaType,
					Source: &errorSource{
						Header: "Content-Type",
					},
				}},
			})
		}
	}
}

func hasContentBody(r *http.Request) bool {
	switch r.Method {
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		return true
	default:
		return false
	}
}

func checkContentType(contentType, mediaType string) bool {
	if contentType == "" {
		// Ignore empty Content-Type headers
		return true
	}

	contentType, params, err := mime.ParseMediaType(contentType)
	if err != nil || contentType != mediaType {
		return false
	}

	if len(params) > 0 {
		// charset with the value utf-8 is the only allowed content type parameter
		if val, ok := params["charset"]; len(params) != 1 || !ok || val != "utf-8" {
			return false
		}
	}

	return true
}

func acceptHeaderForMethod(method string) (string, bool) {
	switch method {
	case http.MethodPatch:
		return "Accept-Patch", true
	case http.MethodPost:
		return "Accept-Post", true
	default:
		return "", false
	}
}

// panicRecover intercepts and logs panics that occur in request handlers.
func panicRecover(log *log.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					internalError("Recovered from panic", errors.Errorf("%v", e), log)(w, r)
				}
			}()

			next.ServeHTTP(w, r)
		}
	}
}
