package main

import (
	"os"

	"github.com/cerfical/muzik/internal/args"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/httpserv/middleware"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	log := log.New().WithContextKey(middleware.RequestID)
	args := args.Parse(os.Args)

	log = log.WithLevel(args.LogLevel)
	server := httpserv.Server{
		Addr:       args.ServerAddr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}

	server.Use(middleware.LogRequest(log))
	server.Use(middleware.AddRequestID)

	log.WithFields("addr", server.Addr).
		Info("Starting up the server")

	defer func() {
		log.Info("Shutting down the server")
		if err := server.Close(); err != nil {
			log.WithError(err).
				Error("Server shutdown failed")
		}
	}()

	if err := server.Run(); err != nil {
		log.WithError(err).
			Error("Server terminated abnormally")
	}
}
