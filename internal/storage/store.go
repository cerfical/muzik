package storage

import (
	"sync"

	"github.com/cerfical/muzik/internal/model"
)

type Store struct {
	data   []*model.Track
	dataMu sync.Mutex
}

func (s *Store) Create(title string) *model.Track {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	track := &model.Track{
		ID: len(s.data),
		Attrs: model.TrackAttrs{
			Title: title,
		},
	}
	s.data = append(s.data, track)
	return track
}

func (s *Store) Get(id int) (*model.Track, bool) {
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

func (s *Store) GetAll() []*model.Track {
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
