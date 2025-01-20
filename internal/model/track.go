package model

type Track struct {
	ID    int    `json:"id,string"`
	Title string `json:"title"`
}

type TrackInfo struct {
	Title string `json:"title"`
}

type TrackStore interface {
	CreateTrack(info *TrackInfo) *Track
	TrackByID(id int) (*Track, bool)
	AllTracks() []*Track
}
