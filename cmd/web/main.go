package main

import (
	"context"
	"net/http"
	"os"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
)

func main() {
	log := log.New()
	config, err := config.Load(os.Args)
	if err != nil {
		log.Fatal("error loading configuration", err)
	}
	log = log.WithLevel(config.Log.Level)

	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	server := httpserv.New(&config.Server, index, log)
	if err := server.Run(context.Background()); err != nil {
		log.Error("the server terminated abnormally", err)
	}
}
