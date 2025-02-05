package model

import (
	"context"
	"io"
)

type Track struct {
	ID    int        `json:"id,string"`
	Attrs TrackAttrs `json:"attributes"`
}

type TrackAttrs struct {
	Title string `json:"title"`
}

type TrackStore interface {
	io.Closer

	CreateTrack(context.Context, *TrackAttrs) (*Track, error)
	TrackByID(context.Context, int) (*Track, error)
	AllTracks(context.Context) ([]Track, error)
}
