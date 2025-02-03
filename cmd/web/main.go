package main

import (
	"net/http"
	"os"

	"github.com/cerfical/muzik/internal/config"
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

	log.WithFields("addr", config.Server.Addr).Info("starting up the WEB server")
	if err := http.ListenAndServe(config.Server.Addr, index); err != nil {
		log.Error("the WEB server terminated abnormally", err)
	}
	log.Info("shutting down the WEB server")
}
