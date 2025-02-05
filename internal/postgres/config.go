package postgres

import "time"

type Config struct {
	Addr string
	Name string

	User     string
	Password string

	Timeout     time.Duration
	IdleTimeout time.Duration
}
