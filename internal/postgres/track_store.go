package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

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
	db.SetConnMaxIdleTime(cfg.IdleTimeout)

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tracks(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL
		)
	`); err != nil {
		return nil, err
	}

	return &TrackStore{db, cfg.Timeout}, nil
}

func makeConnString(cfg *Config) (string, error) {
	host, port, err := net.SplitHostPort(cfg.Addr)
	if cfg.Addr != "" && err != nil {
		return "", err
	}

	c := []struct {
		key, val string
	}{
		{"host", host},
		{"port", port},
		{"user", cfg.User},
		{"password", cfg.Password},
		{"database", cfg.Name},
		{"sslmode", "disable"},
	}

	var options []string
	for _, cc := range c {
		if cc.val == "" {
			continue
		}
		options = append(options, fmt.Sprintf("%v='%v'", cc.key, cc.val))
	}

	connStr := strings.Join(options, " ")
	return connStr, nil
}

type TrackStore struct {
	db      *sql.DB
	timeout time.Duration
}

func (s *TrackStore) CreateTrack(ctx context.Context, attrs *model.TrackAttrs) (*model.Track, error) {
	var id int
	err := s.withTimeout(ctx, func(ctx context.Context) error {
		row := s.db.QueryRowContext(ctx,
			"INSERT INTO tracks(title) VALUES($1) RETURNING id",
			attrs.Title,
		)
		return row.Scan(&id)
	})

	if err != nil {
		return nil, err
	}

	return &model.Track{ID: id, Attrs: *attrs}, nil
}

func (s *TrackStore) TrackByID(ctx context.Context, id int) (*model.Track, error) {
	var track model.Track
	err := s.withTimeout(ctx, func(ctx context.Context) error {
		row := s.db.QueryRowContext(ctx, "SELECT id, title FROM tracks WHERE id=$1", id)
		return row.Scan(&track.ID, &track.Attrs.Title)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &track, nil
}

func (s *TrackStore) AllTracks(ctx context.Context) ([]model.Track, error) {
	var tracks []model.Track
	err := s.withTimeout(ctx, func(ctx context.Context) (err error) {
		rows, err := s.db.QueryContext(ctx, "SELECT id, title FROM tracks")
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		for rows.Next() {
			var track model.Track
			if err = rows.Scan(&track.ID, &track.Attrs.Title); err != nil {
				return err
			}
			tracks = append(tracks, track)
		}
		return rows.Err()
	})

	return tracks, err
}

func (s *TrackStore) withTimeout(ctx context.Context, f func(ctx context.Context) error) error {
	timedCtx := ctx
	if s.timeout > 0 {
		var cancel context.CancelFunc
		timedCtx, cancel = context.WithTimeout(timedCtx, s.timeout)
		defer cancel()
	}

	return f(timedCtx)
}

func (s *TrackStore) Close() error {
	return s.db.Close()
}
