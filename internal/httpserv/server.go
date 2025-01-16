package httpserv

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

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
	serv := http.Server{
		Addr:    s.Addr,
		Handler: s.setupRouter(),
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	shutdownChan := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(sigChan)

		// Make sure the goroutine responsible for graceful shutdown is terminated properly
		close(sigChan)
		<-shutdownChan
	}()

	go func() {
		defer close(shutdownChan)

		// If the server was terminated due to some other error, return immediately
		if _, ok := <-sigChan; !ok {
			return
		}

		// Try to shutdown the server cleanly and if that fails, close the server
		if err := serv.Shutdown(context.Background()); err != nil {
			s.Log.WithError(err).
				Error("graceful shutdown failed")

			if err := serv.Close(); err != nil {
				s.Log.WithError(err).
					Error("server shutdown failed")
			}
		}
	}()

	if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) setupRouter() http.Handler {
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
	return http.HandlerFunc(requestLogger(log, mux))
}

func requestLogger(l *log.Logger, next http.Handler) http.HandlerFunc {
	var requestID atomic.Uint64

	return func(w http.ResponseWriter, r *http.Request) {
		id := requestID.Add(1)
		ctx := context.WithValue(r.Context(), keyRequestID{}, id)

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
