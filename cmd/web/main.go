package main

import (
	"net/http"
	"os"

	"github.com/cerfical/muzik/internal/middleware"
	"github.com/cerfical/muzik/internal/webapp"
)

func main() {
	app := webapp.New(os.Args)

	app.Route("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	app.Use(middleware.LogRequest(app.Log))

	app.Run()
}
