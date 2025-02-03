package main

import (
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

	if err := httpserv.Run(&config.Server, index, log); err != nil {
		log.Error("the server terminated abnormally", err)
	}
}
