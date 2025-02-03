package api

import (
	"net/http"

	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/model"
)

func New(store model.TrackStore, log *log.Logger) http.Handler {
	return setupRouter(store, log)
}
