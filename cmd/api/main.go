package main

import (
	"os"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/httpserv/middleware"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

type Config struct {
	Log struct {
		Level log.Level `mapstructure:"level"`
	} `mapstructure:"log"`
}

func main() {
	log := log.New().WithContextKey(middleware.RequestID)

	config, err := config.Load(os.Args)
	if err != nil {
		log.WithError(err).
			Fatal("error loading config file")
	}

	log = log.WithLevel(config.Log.Level)
	log.WithFields(
		"server.addr", config.Server.Addr,
		"log.level", config.Log.Level,
	).Info("Using config")

	server := httpserv.Server{
		Addr:       config.Server.Addr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}

	server.Use(middleware.LogRequest(log))
	server.Use(middleware.AddRequestID)

	log.Info("Starting up the server")
	if err := server.Run(); err != nil {
		log.WithError(err).
			Fatal("Server terminated abnormally")
	}

	log.Info("Shutting down the server")
	if err := server.Close(); err != nil {
		log.WithError(err).
			Fatal("Server shutdown failed")
	}
}
