package api

import (
	"net/http"

	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

func NewServer(config *httpserv.Config, store model.TrackStore, log *log.Logger) *httpserv.Server {
	h := NewHandler(store, log)
	return httpserv.New(config, h, log)
}

func NewHandler(store model.TrackStore, log *log.Logger) http.Handler {
	return setupRouter(store, log)
}
