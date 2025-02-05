package model

import "io"

type Track struct {
	ID    int    `json:"id,string"`
	Title string `json:"title"`
}

type TrackStore interface {
	io.Closer

	CreateTrack(*Track) error
	TrackByID(int) (*Track, error)
	AllTracks() ([]Track, error)
}
