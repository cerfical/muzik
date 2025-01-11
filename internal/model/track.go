package model

type Track struct {
	ID    int         `json:"id,string"`
	Attrs *TrackAttrs `json:"attributes"`
}

type TrackAttrs struct {
	Title string `json:"title"`
}
