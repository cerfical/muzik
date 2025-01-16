package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New() *Logger {
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
	})

	return &Logger{
		logger: zerolog.New(out).With().
			Timestamp().
			Logger(),
	}
}

type Logger struct {
	logger zerolog.Logger
}

func (l *Logger) Fatal(msg string) {
	l.log(LevelFatal, msg)
	os.Exit(1)
}

func (l *Logger) Error(msg string) {
	l.log(LevelError, msg)
}

func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg)
}

func (l *Logger) log(lvl Level, msg string) {
	l.logger.WithLevel(zerolog.Level(lvl)).Msg(msg)
}

func (l *Logger) WithLevel(lvl Level) *Logger {
	return &Logger{l.logger.Level(zerolog.Level(lvl))}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{l.logger.With().Ctx(ctx).Logger()}
}

func (l *Logger) WithContextKey(key any) *Logger {
	f := func(e *zerolog.Event, _ zerolog.Level, _ string) {
		if v := e.GetCtx().Value(key); v != nil {
			e.Any(fmt.Sprintf("%v", key), v)
		}
	}
	return &Logger{l.logger.Hook(zerolog.HookFunc(f))}
}

func (l *Logger) WithFields(fields ...any) *Logger {
	if len(fields)%2 != 0 {
		panic("expected an even number of arguments")
	}
	return &Logger{l.logger.With().Fields(fields).Logger()}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{l.logger.With().Err(err).Logger()}
}
