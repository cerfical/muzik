package middleware

import (
	"net/http"

	"github.com/cerfical/muzik/internal/log"
)

func LogRequest(l *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ll := l.WithFields(
				"method", r.Method,
				"path", r.URL.Path,
			)
			ll.Info("Incoming request")

			rr := r.WithContext(ll.WithContext(r.Context()))
			next.ServeHTTP(w, rr)
		})
	}
}
