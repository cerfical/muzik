package httpserv

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type Server struct {
	TrackStore      model.TrackStore
	Addr            string
	Log             *log.Logger
	shutdownErrChan chan error
	middleware      []Middleware
}

type Middleware func(http.Handler) http.Handler

func (s *Server) Use(m Middleware) {
	s.middleware = append(s.middleware, m)
}

func (s *Server) Run() error {
	h := s.setupRouter()
	for _, m := range s.middleware {
		h = m(h)
	}

	serv := http.Server{
		Addr:    s.Addr,
		Handler: h,
	}

	// Graceful shutdown
	s.shutdownErrChan = make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		// Make sure the shutdown goroutine is terminated properly
		signal.Stop(sigChan)
		close(sigChan)
	}()

	go func() {
		var err error
		defer func() {
			s.shutdownErrChan <- err
		}()

		// If the server was terminated due to some other error, return immediately
		if _, ok := <-sigChan; !ok {
			return
		}

		// Try to shutdown the server cleanly and if that fails, close the server
		err = serv.Shutdown(context.Background())
		if err != nil {
			serv.Close()
		}
	}()

	if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) setupRouter() http.Handler {
	tracks := api.TrackHandler{
		Store: s.TrackStore,
		Log:   s.Log,
	}

	mux := http.NewServeMux()
	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"GET /api/tracks/{id}", tracks.Get},
		{"GET /api/tracks/{$}", tracks.GetAll},
		{"POST /api/tracks/{$}", tracks.Create},
		{"GET /{$}", index},
	}

	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}
	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func (s *Server) Close() error {
	defer close(s.shutdownErrChan)
	return <-s.shutdownErrChan
}
