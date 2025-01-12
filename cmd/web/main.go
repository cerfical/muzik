package main

import (
	"log"

	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/memstore"
)

func main() {
	server := httpserv.Server{
		Addr:       "127.0.0.1:8080",
		TrackStore: &memstore.TrackStore{},
	}

	if err := server.Run(); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}
