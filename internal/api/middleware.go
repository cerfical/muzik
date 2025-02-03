package api

import (
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/strutil"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// allowMethods populates Allow header with methods defined for the request path on a [mux.Router].
func allowMethods(router *mux.Router, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowedMethods []string
		router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
			path, _ := route.GetPathTemplate()

			matcher := mux.NewRouter()
			matcher.Path(path)

			var match mux.RouteMatch
			if matcher.Match(r, &match) {
				methods, _ := route.GetMethods()
				allowedMethods = append(allowedMethods, methods...)
			}
			return nil
		})

		allowedMethods = strutil.Dedup(allowedMethods)
		w.Header().Set("Allow", strings.Join(allowedMethods, ", "))

		next(w, r)
	})
}

// fillPathParams makes gorilla/mux path parameters available via [http.Request.PathValue].
func fillPathParams(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, val := range mux.Vars(r) {
			r.SetPathValue(key, val)
		}

		next(w, r)
	}
}

// accepts checks Accept header for the presence of the specified media type.
func accepts(mediaType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkAcceptHeader(mediaType, r.Header.Get("Accept")) {
				next.ServeHTTP(w, r)
				return
			}

			details := fmt.Sprintf("the only acceptable media type is '%s'", mediaType)
			encodeError(w, &apiError{
				Title:  "media type is not acceptable",
				Detail: details,
				Status: http.StatusNotAcceptable,
				Source: &errorSource{
					Header: "Accept",
				},
			})
		})
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
func hasContentType(mediaType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if !hasContentBody(r) || checkContentType(contentType, mediaType) {
				next.ServeHTTP(w, r)
				return
			}

			if h, ok := acceptHeaderForMethod(r.Method); ok {
				w.Header().Set(h, mediaType)
			}

			details := fmt.Sprintf("unexpected content type '%s', only '%s' is allowed", contentType, mediaType)
			encodeError(w, &apiError{
				Title:  "media type is unsupported",
				Detail: details,
				Status: http.StatusUnsupportedMediaType,
				Source: &errorSource{
					Header: "Content-Type",
				},
			})
		})
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

// panicRecover intercepts and logs panics that occur in request handlers.
func panicRecover(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					log.Error("recovered from panic", errors.Errorf("%v", e))
					internalError(w, r)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
