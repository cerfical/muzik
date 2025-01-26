package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/cerfical/muzik/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenTrackStore(cfg *Config) (model.TrackStore, error) {
	connStr, err := makeConnString(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tracks(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL
		)
	`); err != nil {
		return nil, err
	}

	return &TrackStore{db}, nil
}

func makeConnString(cfg *Config) (string, error) {
	host, port, err := net.SplitHostPort(cfg.Addr)
	if cfg.Addr != "" && err != nil {
		return "", err
	}

	// Use slice rather than map for deterministic key order
	c := []struct {
		key, val string
	}{
		{"host", host},
		{"port", port},
		{"user", cfg.User},
		{"password", cfg.Password},
		{"database", cfg.Database},
		{"sslmode", "disable"},
	}

	var connStr string
	for _, cc := range c {
		if cc.val == "" {
			continue
		}

		// Separate settings with spaces
		if connStr != "" {
			connStr += " "
		}
		connStr += fmt.Sprintf("%v='%v'", cc.key, cc.val)
	}
	return connStr, nil
}

type TrackStore struct {
	db *sql.DB
}

func (s *TrackStore) CreateTrack(track *model.Track) error {
	row := s.db.QueryRow("INSERT INTO tracks(title) VALUES($1) RETURNING id", track.Title)
	return row.Scan(&track.ID)
}

func (s *TrackStore) TrackByID(id int) (*model.Track, error) {
	row := s.db.QueryRow("SELECT id, title FROM tracks WHERE id=$1", id)

	var track model.Track
	if err := row.Scan(&track.ID, &track.Title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &track, nil
}

func (s *TrackStore) AllTracks() ([]model.Track, error) {
	rows, err := s.db.Query("SELECT id, title FROM tracks")
	if err != nil {
		return nil, err
	}

	// Return the empty collection as a slice of size 0, not as nil
	tracks := make([]model.Track, 0)

	for rows.Next() {
		var track model.Track
		err = rows.Scan(&track.ID, &track.Title)
		if err != nil {
			break
		}
		tracks = append(tracks, track)
	}

	// Check for close errors
	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	// Check for scan errors
	if err != nil {
		return nil, err
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s *TrackStore) Close() error {
	return s.db.Close()
}
