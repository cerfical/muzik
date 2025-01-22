package middleware

import (
	"mime"
	"net/http"
	"strings"
)

func HasContentType(mediaType string) func(http.Handler) http.Handler {
	mediaType = strings.ToLower(mediaType)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasContentBody(r) && !hasContentType(r, mediaType) {
				w.Header().Set(acceptHeaderForMethod(r.Method), mediaType)
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasContentType(r *http.Request, mediaType string) bool {
	contentType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != mediaType || len(params) != 0 {
		return false
	}
	return true
}

func hasContentBody(r *http.Request) bool {
	switch r.Method {
	case http.MethodPatch, http.MethodPost:
		return true
	default:
		return false
	}
}

func acceptHeaderForMethod(method string) string {
	switch method {
	case http.MethodPatch:
		return "Accept-Patch"
	case http.MethodPost:
		return "Accept-Post"
	default:
		return ""
	}
}
