package main

import (
	"log"

	"github.com/cerfical/muzik/internal/httpserv"
)

func main() {
	server := httpserv.New()
	if err := server.Run("127.0.0.1:8080"); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}
