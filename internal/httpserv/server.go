package httpserv

import (
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv/api/tracks"
	"github.com/cerfical/muzik/internal/storage"
)

func New() *Server {
	server := &Server{mux: http.NewServeMux()}
	routes := []struct {
		path    string
		handler http.Handler
	}{
		{"GET /api/tracks/{id}", tracks.Get(&server.store)},
		{"GET /api/tracks/", tracks.List(&server.store)},
		{"GET /{$}", http.HandlerFunc(server.index)},
	}

	for _, r := range routes {
		server.mux.Handle(r.path, r.handler)
	}
	return server
}

type Server struct {
	mux   *http.ServeMux
	store storage.Store
}

func (s *Server) index(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "static/index.html")
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
