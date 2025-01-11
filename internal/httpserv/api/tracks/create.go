package tracks

import (
	"net/http"

	"github.com/cerfical/muzik/internal/storage"
)

type trackCreator struct {
	store *storage.Store
}

func Create(store *storage.Store) http.Handler {
	return &trackCreator{store}
}

func (h *trackCreator) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	attrs, err := readTrackAttrs(req.Body)
	if err != nil {
		writeErrorResponse(wr, http.StatusBadRequest, err.Error())
		return
	}

	track := h.store.Create(attrs)
	writeTrack(wr, track)
}
