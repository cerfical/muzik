package tracks

import (
	"encoding/json"
	"io"

	"github.com/cerfical/muzik/internal/model"
)

func readTrackAttrs(r io.Reader) (*model.TrackAttrs, error) {
	var req struct {
		Data *model.TrackAttrs `json:"data"`
	}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	return req.Data, nil
}
