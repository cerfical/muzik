package main

import (
	"os"

	"github.com/cerfical/muzik/internal/api"
	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/postgres"
)

func main() {
	log := log.New()
	config, err := config.Load(os.Args)
	if err != nil {
		log.Fatal("error loading configuration", err)
	}
	log = log.WithLevel(config.Log.Level)

	log.WithFields(
		"addr", config.DB.Addr,
		"user", config.DB.User,
		"name", config.DB.Name,
	).Info("opening the database")

	store, err := postgres.OpenTrackStore(&config.DB)
	if err != nil {
		log.Fatal("failed to open the database", err)
	}

	defer func() {
		log.Info("closing the database")
		if err := store.Close(); err != nil {
			log.Error("failed to close the database", err)
		}
	}()

	server := api.NewServer(store, log)
	if err := server.Run(config.Server.Addr); err != nil {
		log.Error("the API server terminated abnormally", err)
	}
}
