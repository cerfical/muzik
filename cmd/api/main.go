package main

import (
	"os"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/httpserv/middleware"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/postgres"
)

func main() {
	log := log.New().WithContextKey(middleware.RequestID)

	config, err := config.Load(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Error loading config file")
	}

	log = log.WithLevel(config.Log.Level)
	log.WithFields(
		"server.addr", config.Server.Addr,
		"storage.addr", config.Storage.Addr,
		"storage.user", config.Storage.User,
		"storage.database", config.Storage.Database,
		"log.level", config.Log.Level,
	).Info("Using config")

	log.Info("Opening the database")
	store, err := postgres.OpenTrackStore(&config.Storage)
	if err != nil {
		log.WithError(err).Fatal("Failed to open the database")
	}

	server := httpserv.Server{
		Addr:       config.Server.Addr,
		TrackStore: store,
		Log:        log,
	}

	server.Use(middleware.LogRequest(log))
	server.Use(middleware.AddRequestID)

	log.Info("Starting up the server")
	if err := server.Run(); err != nil {
		log.WithError(err).Fatal("Server terminated abnormally")
	}

	log.Info("Shutting down the server")
	if err := server.Close(); err != nil {
		log.WithError(err).Fatal("Server shutdown failed")
	}

	log.Info("Closing the database")
	if err := store.Close(); err != nil {
		log.WithError(err).Fatal("Failed to close the database")
	}
}
