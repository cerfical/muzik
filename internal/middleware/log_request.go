package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cerfical/muzik/internal/log"
)

func LogRequest(l *log.Logger) func(http.Handler) http.Handler {
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
