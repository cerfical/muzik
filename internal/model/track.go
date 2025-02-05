package model

import (
	"context"
	"io"
)

type Track struct {
	ID    int    `json:"id,string"`
	Title string `json:"title"`
}

type TrackStore interface {
	io.Closer

	CreateTrack(context.Context, *Track) error
	TrackByID(context.Context, int) (*Track, error)
	AllTracks(context.Context) ([]Track, error)
}
