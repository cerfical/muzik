package main

import (
	"log"

	"github.com/cerfical/muzik/internal/serv"
)

func main() {
	server, err := serv.New()
	if err != nil {
		log.Fatalf("server startup: %v", err)
	}

	if err := server.Run("127.0.0.1:8080"); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}
