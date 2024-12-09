package storage

// Track represents information about a single music track
type Track struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}
