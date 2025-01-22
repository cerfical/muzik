package httpserv

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Addr    string
	Handler http.Handler

	shutdownErrChan chan error
}

func (s *Server) Run() error {
	// Graceful shutdown
	s.shutdownErrChan = make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		// Make sure the shutdown goroutine is terminated properly
		signal.Stop(sigChan)
		close(sigChan)
	}()

	serv := http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}

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
			s.Close()
		}
	}()

	if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	return <-s.shutdownErrChan
}
