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
	log := log.New()
	config, err := config.Load(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Error loading config file")
	}

	log = log.WithLevel(config.Log.Level)
	serv := httpserv.Server{
		Addr: config.Server.Addr,
		Log:  log,
	}

	// Setup middleware
	serv.Use(middleware.LogRequest(log))

	// Setup routes
	serv.Route("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	serv.Run()
}
