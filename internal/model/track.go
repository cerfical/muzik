package model

type Track struct {
	ID    int    `json:"id,string"`
	Title string `json:"title"`
}

type TrackStore interface {
	CreateTrack(*Track) error
	TrackByID(int) (*Track, error)
	AllTracks() ([]Track, error)
	Close() error
}
