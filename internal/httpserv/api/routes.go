package api

import (
	"net/http"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"

	"github.com/gorilla/mux"
)

func setupRouter(store model.TrackStore, log *log.Logger) http.Handler {
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = allowMethods(router, methodNotAllowed)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	router.Use(hasContentType(encodeMediaType))
	router.Use(accepts(encodeMediaType))
	router.Use(panicRecover(log))

	tracks := tracksHandler{store, log}
	router.HandleFunc("/api/tracks/{id}", fillPathParams(tracks.get)).Methods("GET")
	router.HandleFunc("/api/tracks/{id}", fillPathParams(tracks.delete)).Methods("DELETE")
	router.HandleFunc("/api/tracks/", tracks.getAll).Methods("GET")
	router.HandleFunc("/api/tracks/", tracks.create).Methods("POST")

	return router
}
