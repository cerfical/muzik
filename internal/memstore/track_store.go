package memstore

import (
	"sync"

	"github.com/cerfical/muzik/internal/model"
)

type TrackStore struct {
	data   []*model.Track
	dataMu sync.Mutex
}

func (s *TrackStore) CreateTrack(info *model.TrackInfo) *model.Track {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	t := &model.Track{
		ID:   len(s.data),
		Info: info,
	}
	s.data = append(s.data, t)
	return t
}

func (s *TrackStore) TrackByID(id int) (*model.Track, bool) {
	if id < 0 {
		return nil, false
	}

	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	if id >= len(s.data) {
		return nil, false
	}
	return s.data[id], true
}

func (s *TrackStore) AllTracks() []*model.Track {
	d := func() []*model.Track {
		s.dataMu.Lock()
		defer s.dataMu.Unlock()
		return s.data
	}()

	if d == nil {
		return []*model.Track{}
	}
	return d
}
