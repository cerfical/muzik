package middleware

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/cerfical/muzik/internal/api/errors"
)

func HasContentType(mediaType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if !hasContentBody(r) || checkContentType(contentType, mediaType) {
				next(w, r)
				return
			}

			if h, ok := acceptHeaderForMethod(r.Method); ok {
				w.Header().Set(h, mediaType)
			}

			details := fmt.Sprintf("unexpected content type '%s', only '%s' is allowed", contentType, mediaType)
			e := errors.Error{
				Title:  "media type is unsupported",
				Detail: details,
				Status: http.StatusUnsupportedMediaType,
			}
			e.Write(w)
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
	if err != nil || contentType != mediaType || len(params) > 0 {
		return false
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
