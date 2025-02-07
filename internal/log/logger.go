package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func New() *Logger {
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
	})

	return &Logger{
		logger: zerolog.New(out).With().
			Stack().
			Timestamp().
			Logger(),
	}
}

func init() {
	// Free the "time" field name for request timing
	zerolog.TimestampFieldName = "timestamp"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

type Logger struct {
	logger zerolog.Logger
}

func (l *Logger) Fatal(msg string, err error) {
	l.log(LevelFatal, msg, err)
	os.Exit(1)
}

func (l *Logger) Error(msg string, err error) {
	l.log(LevelError, msg, err)
}

func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg, nil)
}

func (l *Logger) log(lvl Level, msg string, err error) {
	if l == nil {
		return
	}

	logEv := l.logger.WithLevel(zerolog.Level(lvl))
	if err != nil {
		logEv = logEv.Err(err)
	}
	logEv.Msg(msg)
}

func (l *Logger) WithLevel(lvl Level) *Logger {
	if l == nil {
		return l
	}

	return &Logger{l.logger.Level(zerolog.Level(lvl))}
}

func (l *Logger) WithFields(fields ...any) *Logger {
	if len(fields)%2 != 0 {
		panic("expected an even number of arguments")
	}

	if l == nil {
		return l
	}

	return &Logger{l.logger.With().Fields(fields).Logger()}
}
