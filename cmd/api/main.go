package main

import (
	"context"
	"os"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/postgres"
)

func main() {
	config := config.MustLoad(os.Args)

	log := log.New().WithLevel(config.Log.Level)
	log.WithFields(
		"addr", config.DB.Addr,
		"user", config.DB.User,
		"name", config.DB.Name,
	).Info("Opening the database")

	store, err := postgres.OpenTrackStore(&config.DB)
	if err != nil {
		log.Fatal("Failed to open the database", err)
	}

	defer func() {
		log.Info("Closing the database")
		if err := store.Close(); err != nil {
			log.Error("Failed to close the database", err)
		}
	}()

	server := api.NewServer(&config.Server, store, log)
	if err := server.Run(context.Background()); err != nil {
		log.Error("The server has terminated abnormally", err)
	}
}
