package middleware

import (
	"net/http"
	"sync/atomic"

	"github.com/cerfical/muzik/internal/log"
)

func LogRequest(l *log.Logger) func(http.Handler) http.Handler {
	var requestCount atomic.Uint64

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := requestCount.Add(1)

			ll := l.WithFields("request_id", requestID)
			ll.WithFields(
				"method", r.Method,
				"path", r.URL.Path,
			).Info("Incoming request")

			rr := r.WithContext(ll.WithContext(r.Context()))
			next.ServeHTTP(w, rr)
		})
	}
}
