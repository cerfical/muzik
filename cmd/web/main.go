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
	config := config.MustLoad(os.Args)
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	log := log.New().WithLevel(config.Log.Level)
	server := httpserv.New(&config.Server, index, log)
	if err := server.Run(context.Background()); err != nil {
		log.Error("The server has terminated abnormally", err)
	}
}
