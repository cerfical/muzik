package middleware

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/cerfical/muzik/internal/log"
)

func LogRequest(l *log.Logger) func(http.Handler) http.Handler {
	var requestCount atomic.Uint64

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := requestCount.Add(1)

			// Log the start of request processing
			ll := l.WithFields("request_id", requestID)
			ll.WithFields(
				"method", r.Method,
				"path", r.URL.Path,
			).Info("Incoming request")

			rr := r.WithContext(ll.WithContext(r.Context()))
			ww := loggingResponseWriter{ResponseWriter: w}

			// Time the request
			startTime := time.Now()
			next.ServeHTTP(&ww, rr)
			elapsed := time.Since(startTime)

			// Log the end of request processing
			statusLine := fmt.Sprintf("%d %s", ww.StatusCode, http.StatusText(ww.StatusCode))
			ll.WithFields(
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
