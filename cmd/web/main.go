package main

import (
	"context"

	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	servAddr := "127.0.0.1:8080"
	log := log.New()

	ctx := context.Background()
	server := httpserv.Server{
		Addr:       servAddr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}
	log.WithString("addr", servAddr).Info(ctx, "server startup")

	if err := server.Run(); err != nil {
		log.WithError(err).Fatal(ctx, "server shutdown")
	}
}
