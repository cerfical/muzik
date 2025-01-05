package stor

import "sync"

// Store stores tracks in some kind of storage, such as a database
type Store struct {
	data   []*Track
	dataMu sync.Mutex
}

// Create creates a new named track in the store
func (s *Store) Create(title string) *Track {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	track := &Track{
		ID:    len(s.data),
		Title: title,
	}
	s.data = append(s.data, track)
	return track
}

// Get retrieves a track from the store by its ID
func (s *Store) Get(id int) (*Track, bool) {
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

// GetAll retrieves all tracks from the store
func (s *Store) GetAll() []*Track {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	return s.data
}
