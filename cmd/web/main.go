package main

import (
	"os"

	"github.com/cerfical/muzik/internal/args"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	log := log.New()
	args := args.Parse(os.Args)

	log = log.WithLevel(args.LogLevel)
	server := httpserv.Server{
		Addr:       args.ServerAddr,
		TrackStore: &memstore.TrackStore{},
		Log:        log,
	}

	log.WithString("addr", args.ServerAddr).
		Info("server startup")

	if err := server.Run(); err != nil {
		log.WithError(err).
			Fatal("server shutdown")
	}
}
