package webapp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/strutil"
	"github.com/gorilla/mux"
)

// allowMethods collects defined methods from a [mux.Router] and sets the Allow header accordingly.
func allowMethods(router *mux.Router, next http.Handler) http.Handler {
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

		next.ServeHTTP(w, r)
	})
}

// fillPathParams makes gorilla/mux path parameters available via [http.Request.PathValue].
func fillPathParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range mux.Vars(r) {
			r.SetPathValue(key, val)
		}

		next.ServeHTTP(w, r)
	})
}

// logRequest logs all incoming requests to the specified [log.Logger].
func logRequest(l *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := loggingResponseWriter{ResponseWriter: w}

			// Time the request
			startTime := time.Now()
			next.ServeHTTP(&ww, r)
			elapsed := time.Since(startTime)

			// Log the end of request processing
			statusLine := fmt.Sprintf("%d %s", ww.StatusCode, http.StatusText(ww.StatusCode))
			l.WithFields(
				"method", r.Method,
				"path", r.URL.Path,
				"status", statusLine,
				"time", elapsed.String(),
			).Info("Request complete")
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter

	StatusCode int
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
func (w *loggingResponseWriter) Write(buf []byte) (int, error) {
	return w.ResponseWriter.Write(buf)
}
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
