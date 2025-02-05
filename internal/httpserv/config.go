package httpserv

import "time"

type Config struct {
	Addr string

	Timeout     time.Duration
	IdleTimeout time.Duration
}
