package main

import (
	"net/http"
	"os"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/middleware"
)

func main() {
	app := NewApp()
	app.Run()
}

type App struct {
	config *config.Config
	server *httpserv.Server
	log    *log.Logger
}

func NewApp() *App {
	var a App

	a.log = log.New().WithContextKey(middleware.RequestID)
	a.config = loadConfig(a.log)

	a.log = a.log.WithLevel(a.config.Log.Level)
	a.server = newServer(a.config, a.log)

	return &a
}

func loadConfig(l *log.Logger) *config.Config {
	c, err := config.Load(os.Args)
	if err != nil {
		l.WithError(err).Fatal("Error loading config file")
	}
	return c
}

func newServer(c *config.Config, l *log.Logger) *httpserv.Server {
	serv := httpserv.Server{
		Addr: c.Server.Addr,
		Log:  l,
	}

	setupMiddleware(&serv, l)
	setupRoutes(&serv)

	return &serv
}

func setupRoutes(serv *httpserv.Server) {
	serv.Route("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
}

func setupMiddleware(serv *httpserv.Server, l *log.Logger) {
	serv.Use(middleware.LogRequest(l))
	serv.Use(middleware.AddRequestID)
}

func (a *App) Run() {
	// Dump the app configuration being used
	a.log.WithFields(
		"server.addr", a.config.Server.Addr,
		"log.level", a.config.Log.Level,
	).Info("Using config")

	a.log.Info("Starting up the server")
	if err := a.server.Run(); err != nil {
		a.log.WithError(err).Fatal("Server terminated abnormally")
	}

	a.closeServer()
}

func (a *App) closeServer() {
	a.log.Info("Shutting down the server")
	if err := a.server.Close(); err != nil {
		a.log.WithError(err).Fatal("Server shutdown failed")
	}
}
