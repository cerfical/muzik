package httpserv

import (
	"fmt"
	"net/http"

	"github.com/cerfical/muzik/internal/storage"
)

// New creates and initializes a new server instance
func New() (serv *Server, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	server := &Server{mux: http.NewServeMux()}
	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"GET /api/tracks/{id}", server.trackByID},
		{"GET /api/tracks/", server.listTracks},
		{"GET /{$}", server.index},
	}

	for _, r := range routes {
		server.mux.HandleFunc(r.path, r.handler)
	}
	return server, nil
}

// Server is an instance of the API server
type Server struct {
	mux   *http.ServeMux
	store storage.Store
}

// Run starts the server at the specified address
func (s *Server) Run(servAddr string) error {
	return http.ListenAndServe(servAddr, s.mux)
}
