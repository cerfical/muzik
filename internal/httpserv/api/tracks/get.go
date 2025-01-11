package tracks

import (
	"net/http"
	"strconv"

	"github.com/cerfical/muzik/internal/storage"
)

func Get(store *storage.Store) http.Handler {
	return &trackGetter{store}
}

type trackGetter struct {
	store *storage.Store
}

func (h *trackGetter) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	trackID, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		writeErrorResponse(wr, http.StatusNotFound, "track not found")
		return
	}

	track, ok := h.store.Get(trackID)
	if !ok {
		writeErrorResponse(wr, http.StatusNotFound, "track not found")
		return
	}
	writeTrack(wr, track)
}
