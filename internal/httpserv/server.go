package httpserv

import (
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv/api"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

type Server struct {
	TrackStore model.TrackStore
	Addr       string
	Log        *log.Logger
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	tracks := api.TrackHandler{Store: s.TrackStore}
	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"GET /api/tracks/{id}", tracks.Get},
		{"GET /api/tracks/{$}", tracks.GetAll},
		{"POST /api/tracks/{$}", tracks.Create},
		{"GET /{$}", s.index},
	}

	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}

	return http.ListenAndServe(s.Addr, mux)
}

func (s *Server) index(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "static/index.html")
}
