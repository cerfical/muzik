package main

import (
	"os"

	"github.com/cerfical/muzik/internal/api"
	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/middleware"
	"github.com/cerfical/muzik/internal/postgres"
)

func main() {
	log := log.New()
	config, err := config.Load(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Error loading config file")
	}

	log = log.WithLevel(config.Log.Level)
	log.WithFields(
		"stor.addr", config.Storage.Addr,
		"stor.user", config.Storage.User,
		"stor.db", config.Storage.Database,
	).Info("Opening the database")

	// Database configuration
	store, err := postgres.OpenTrackStore(&config.Storage)
	if err != nil {
		log.WithError(err).Fatal("Failed to open the database")
	}

	defer func() {
		log.Info("Closing the database")
		if err := store.Close(); err != nil {
			log.WithError(err).Fatal("Failed to close the database")
		}
	}()

	// Server configuration
	serv := httpserv.Server{
		Addr: config.Server.Addr,
		Log:  log,
	}

	// Setup middleware
	serv.Use(middleware.LogRequest(log))

	// Setup routes
	tracks := api.TrackHandler{Store: store}
	serv.Route("GET /api/tracks/{id}", tracks.Get)
	serv.Route("GET /api/tracks/{$}", tracks.GetAll)
	serv.Route("POST /api/tracks/{$}", tracks.Create)

	serv.Run()
}
