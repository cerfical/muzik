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
	log := log.New(&config.Log)

	server := httpserv.New(&config.Server, http.FileServer(http.Dir("static")), log)
	if err := server.Run(context.Background()); err != nil {
		log.Error("The server has terminated abnormally", err)
	}
}
