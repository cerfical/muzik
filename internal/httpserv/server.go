package httpserv

import (
	"context"
	"errors"
	stdlog "log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cerfical/muzik/internal/log"
)

func New(config *Config, h http.Handler, log *log.Logger) *Server {
	return &Server{
		serv: http.Server{
			Addr: config.Addr,

			// Log requests before any routing logic applies
			Handler:  logRequest(log)(h),
			ErrorLog: stdlog.New(&httpErrorLog{log}, "", 0),

			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
			IdleTimeout:  config.IdleTimeout,
		},
		log:     log,
		timeout: config.Timeout,
	}
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

	w.Error("HTTP server error", errors.New(string(p)))
	return n, nil
}

type Server struct {
	serv    http.Server
	log     *log.Logger
	timeout time.Duration
}

func (s *Server) Run(ctx context.Context) error {
	s.log.WithFields("addr", s.serv.Addr).Info("Starting up the server")

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error)
	go func() {
		if err := s.serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-sigCtx.Done():
		// The server stopped due to a system signal, perform graceful shutdown
	case err := <-errChan:
		// The server was terminated abnormally
		return err
	}

	timedCtx := ctx
	if s.timeout > 0 {
		var cancel context.CancelFunc
		timedCtx, cancel = context.WithTimeout(timedCtx, s.timeout)
		defer cancel()
	}

	// Try to shutdown the server cleanly and if that fails, close the server
	s.log.Info("Shutting down the server")
	if err := s.serv.Shutdown(timedCtx); err != nil {
		s.log.Error("Failed to shut down the server", err)
		s.serv.Close()
		return err
	}

	return nil
}
