package api

import (
	stdlog "log"

	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

func NewServer(store model.TrackStore, log *log.Logger) *Server {
	var s Server
	s.store = store
	s.log = log
	return &s
}

type Server struct {
	store model.TrackStore
	log   *log.Logger
}

func (s *Server) Run(addr string) error {
	s.log.WithFields("addr", addr).Info("starting up the API server")

	router := setupRouter(s.store, s.log)
	serv := http.Server{
		Addr: addr,
		// Log requests before any routing logic applies
		Handler:  logRequest(s.log)(router),
		ErrorLog: stdlog.New(&httpErrorLog{s.log}, "", 0),
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		// The server stopped due to a system signal, perform graceful shutdown
	case err := <-errChan:
		// The server was terminated abnormally
		return err
	}

	// Try to shutdown the server cleanly and if that fails, close the server
	s.log.Info("shutting down the API server")
	if err := serv.Shutdown(context.Background()); err != nil {
		s.log.Error("error shutting down the API server", err)
		serv.Close()
		return err
	}

	return nil
}

type httpErrorLog struct {
	*log.Logger
}

func (w *httpErrorLog) Write(p []byte) (int, error) {
	// Trim carriage return produced by stdlog
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
		n--
	}

	w.Error("error serving HTTP", errors.New(string(p)))
	return n, nil
}
