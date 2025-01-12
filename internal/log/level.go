package log

import (
	"errors"

	"github.com/rs/zerolog"
)

const (
	FatalLevel = Level(zerolog.FatalLevel)
	ErrorLevel = Level(zerolog.ErrorLevel)
	InfoLevel  = Level(zerolog.InfoLevel)
	Disabled   = Level(zerolog.Disabled)
)

type Level zerolog.Level

func (l *Level) UnmarshalText(text []byte) error {
	switch text := string(text); text {
	case "fatal":
		*l = FatalLevel
	case "error":
		*l = ErrorLevel
	case "info":
		*l = InfoLevel
	case "none":
		*l = Disabled
	default:
		return errors.New("unknown log level")
	}
	return nil
}

func (l Level) MarshalText() ([]byte, error) {
	var text string
	switch l {
	case FatalLevel:
		text = "fatal"
	case ErrorLevel:
		text = "error"
	case InfoLevel:
		text = "info"
	case Disabled:
		text = "none"
	default:
		return nil, errors.New("unknown log level")
	}
	return []byte(text), nil
}
