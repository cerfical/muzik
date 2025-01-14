package main

import (
	"context"
	"os"

	"github.com/cerfical/muzik/internal/args"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	log := log.New()
	ctx := context.Background()
	args := args.Parse(os.Args)

	log = log.WithLevel(args.LogLevel)
	server := httpserv.Server{
		Addr:       args.ServerAddr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}

	log.WithString("addr", args.ServerAddr).
		Info(ctx, "server startup")

	if err := server.Run(); err != nil {
		log.WithError(err).
			Fatal(ctx, "server shutdown")
	}
}
