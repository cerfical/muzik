package log

import (
	"errors"

	"github.com/rs/zerolog"
)

const (
	LevelNone  = Level(zerolog.Disabled)
	LevelFatal = Level(zerolog.FatalLevel)
	LevelError = Level(zerolog.ErrorLevel)
	LevelInfo  = Level(zerolog.InfoLevel)
)

type Config struct {
	Level Level
}

type Level zerolog.Level

func (l *Level) UnmarshalText(text []byte) error {
	switch text := string(text); text {
	case "fatal":
		*l = LevelFatal
	case "error":
		*l = LevelError
	case "info":
		*l = LevelInfo
	case "none":
		*l = LevelNone
	default:
		return errors.New("unknown log level")
	}
	return nil
}

func (l Level) MarshalText() ([]byte, error) {
	var text string
	switch l {
	case LevelFatal:
		text = "fatal"
	case LevelError:
		text = "error"
	case LevelInfo:
		text = "info"
	case LevelNone:
		text = "none"
	default:
		return nil, errors.New("unknown log level")
	}
	return []byte(text), nil
}
