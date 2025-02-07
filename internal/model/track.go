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
	GetTrack(context.Context, int) (*Track, error)
	GetTracks(context.Context) ([]Track, error)
	DeleteTrack(context.Context, int) error
}
