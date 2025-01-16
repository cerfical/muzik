package httpserv

import (
	"context"
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type Server struct {
	TrackStore model.TrackStore
	Addr       string
	Log        *log.Logger
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	log := s.Log.WithContextValue(keyRequestID{})
	tracks := api.TrackHandler{
		Store: s.TrackStore,
		Log:   log,
	}

	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"GET /api/tracks/{id}", tracks.Get},
		{"GET /api/tracks/{$}", tracks.GetAll},
		{"POST /api/tracks/{$}", tracks.Create},
		{"GET /{$}", s.index},
	}

	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}

	return http.ListenAndServe(s.Addr, http.HandlerFunc(requestLogger(log, mux)))
}

func requestLogger(l *log.Logger, next http.Handler) http.HandlerFunc {
	requestID := 0

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), keyRequestID{}, requestID)
		requestID++

		l.WithStrings(
			"method", r.Method,
			"path", r.URL.Path,
		).
			WithContext(ctx).
			Info("incoming request")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type keyRequestID struct{}

func (keyRequestID) String() string {
	return "request_id"
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
