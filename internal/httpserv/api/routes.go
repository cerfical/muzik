package api

import (
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv/router"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

func setupRoutes(store model.TrackStore, log *log.Logger) http.Handler {
	tracks := tracksHandler{store, log}
	router := router.New().
		Routes("/api/tracks/{id}", []router.Endpoint{
			{"GET", tracks.get},
			{"DELETE", tracks.delete},
		}).
		Routes("/api/tracks/", []router.Endpoint{
			{"POST", tracks.create},
			{"GET", tracks.getAll},
		}).
		Use(hasContentType(encodeMediaType)).
		Use(accepts(encodeMediaType)).
		Use(panicRecover(log))

	return router
}
