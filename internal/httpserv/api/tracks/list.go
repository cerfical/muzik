package tracks

import (
	"net/http"

	"github.com/cerfical/muzik/internal/storage"
)

func List(store *storage.Store) http.Handler {
	return &trackLister{store}
}

type trackLister struct {
	store *storage.Store
}

func (h *trackLister) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	writeTracks(wr, h.store.GetAll())
}
