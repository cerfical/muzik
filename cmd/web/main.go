package main

import (
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	servAddr := "127.0.0.1:8080"
	log := log.New()

	server := httpserv.Server{
		Addr:       servAddr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}
	log.WithString("addr", servAddr).Info("server startup")

	if err := server.Run(); err != nil {
		log.WithError(err).Fatal("server shutdown")
	}
}
