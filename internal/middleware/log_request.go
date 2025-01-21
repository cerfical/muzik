package middleware

import (
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
)

// LogRequest records incoming requests to a log.
func LogRequest(l *log.Logger) httpserv.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.WithContext(r.Context()).
				WithFields(
					"method", r.Method,
					"path", r.URL.Path,
				).
				Info("Incoming request")

			next.ServeHTTP(w, r)
		})
	}
}
