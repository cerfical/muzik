package main

import (
	"os"

	"github.com/cerfical/muzik/internal/handlers"
	"github.com/cerfical/muzik/internal/middleware"
	"github.com/cerfical/muzik/internal/postgres"
	"github.com/cerfical/muzik/internal/webapp"
)

func main() {
	app := webapp.New(os.Args)

	app.Log.WithFields(
		"addr", app.Config.Storage.Addr,
		"user", app.Config.Storage.User,
		"db_name", app.Config.Storage.Database,
	).Info("Opening the database")

	// Database configuration
	store, err := postgres.OpenTrackStore(&app.Config.Storage)
	if err != nil {
		app.Log.WithError(err).Fatal("Failed to open the database")
	}

	defer func() {
		app.Log.Info("Closing the database")
		if err := store.Close(); err != nil {
			app.Log.WithError(err).Error("Failed to close the database")
		}
	}()

	// Setup routes
	tracks := handlers.Tracks{Store: store}
	app.Route("GET /api/tracks/{id}", tracks.Get)
	app.Route("GET /api/tracks/{$}", tracks.GetAll)
	app.Route("POST /api/tracks/{$}", tracks.Create)

	app.Use(middleware.HasContentType("application/json"))
	app.Use(middleware.LogRequest(app.Log))

	app.Run()
}
