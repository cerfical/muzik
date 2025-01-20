package model

type Track struct {
	ID    int    `json:"id,string"`
	Title string `json:"title"`
}

type TrackInfo struct {
	Title string `json:"title"`
}

type TrackStore interface {
	CreateTrack(attrs *TrackInfo) (int, error)
	TrackByID(id int) (*Track, error)
	AllTracks() ([]*Track, error)
	Close() error
}
