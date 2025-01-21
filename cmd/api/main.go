package main

import (
	"os"

	"github.com/cerfical/muzik/internal/api"
	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/middleware"
	"github.com/cerfical/muzik/internal/model"
	"github.com/cerfical/muzik/internal/postgres"
)

func main() {
	app := NewApp()
	app.Run()
}

type App struct {
	config *config.Config
	server *httpserv.Server
	store  model.TrackStore
	log    *log.Logger
}

func NewApp() *App {
	var a App

	a.log = log.New().WithContextKey(middleware.RequestID)
	a.config = loadConfig(a.log)

	a.log = a.log.WithLevel(a.config.Log.Level)
	a.store = newStore(a.config, a.log)
	a.server = newServer(a.config, a.store, a.log)

	return &a
}

func loadConfig(l *log.Logger) *config.Config {
	c, err := config.Load(os.Args)
	if err != nil {
		l.WithError(err).Fatal("Error loading config file")
	}
	return c
}

func newStore(c *config.Config, l *log.Logger) model.TrackStore {
	l.Info("Opening the database")

	s, err := postgres.OpenTrackStore(&c.Storage)
	if err != nil {
		l.WithError(err).Fatal("Failed to open the database")
	}
	return s
}

func newServer(c *config.Config, s model.TrackStore, l *log.Logger) *httpserv.Server {
	serv := httpserv.Server{
		Addr: c.Server.Addr,
		Log:  l,
	}

	setupMiddleware(&serv, l)
	setupRoutes(&serv, s, l)

	return &serv
}

func setupRoutes(serv *httpserv.Server, s model.TrackStore, l *log.Logger) {
	tracks := api.TrackHandler{
		Store: s,
		Log:   l,
	}

	serv.Route("GET /api/tracks/{id}", tracks.Get)
	serv.Route("GET /api/tracks/{$}", tracks.GetAll)
	serv.Route("POST /api/tracks/{$}", tracks.Create)
}

func setupMiddleware(serv *httpserv.Server, l *log.Logger) {
	serv.Use(middleware.LogRequest(l))
	serv.Use(middleware.AddRequestID)
}

func (a *App) Run() {
	// Dump the app configuration being used
	a.log.WithFields(
		"server.addr", a.config.Server.Addr,
		"storage.addr", a.config.Storage.Addr,
		"storage.user", a.config.Storage.User,
		"storage.database", a.config.Storage.Database,
		"log.level", a.config.Log.Level,
	).Info("Using config")

	a.log.Info("Starting up the server")
	if err := a.server.Run(); err != nil {
		a.log.WithError(err).Fatal("Server terminated abnormally")
	}

	a.cleanup()
}

func (a *App) cleanup() {
	a.closeStore()
	a.closeServer()
}

func (a *App) closeStore() {
	a.log.Info("Closing the database")
	if err := a.store.Close(); err != nil {
		a.log.WithError(err).Fatal("Failed to close the database")
	}
}

func (a *App) closeServer() {
	a.log.Info("Shutting down the server")
	if err := a.server.Close(); err != nil {
		a.log.WithError(err).Fatal("Server shutdown failed")
	}
}
