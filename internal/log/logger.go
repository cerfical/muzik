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

func (l *Logger) Fatal(ctx context.Context, msg string) {
	l.log(ctx, LevelFatal, msg)
	os.Exit(1)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.log(ctx, LevelError, msg)
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.log(ctx, LevelInfo, msg)
}

func (l *Logger) log(ctx context.Context, lvl Level, msg string) {
	l.logger.WithLevel(zerolog.Level(lvl)).Ctx(ctx).Msg(msg)
}

func (l *Logger) WithLevel(lvl Level) *Logger {
	return &Logger{l.logger.Level(zerolog.Level(lvl))}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{l.logger.With().Err(err).Logger()}
}

func (l *Logger) WithString(key, val string) *Logger {
	return &Logger{l.logger.With().Str(key, val).Logger()}
}

func (l *Logger) WithInt(key string, val int) *Logger {
	return &Logger{l.logger.With().Int(key, val).Logger()}
}

func (l *Logger) WithContextValue(ctxKey any) *Logger {
	f := func(e *zerolog.Event, _ zerolog.Level, _ string) {
		if v := e.GetCtx().Value(ctxKey); v != nil {
			e.Any(fmt.Sprintf("%v", ctxKey), v)
		}
	}
	return &Logger{l.logger.Hook(zerolog.HookFunc(f))}
}

func (l *Logger) WithStrings(attrs ...string) *Logger {
	if len(attrs)%2 != 0 {
		panic("expected an even number of arguments")
	}

	ctxLogger := l
	for i := 0; i < len(attrs); i += 2 {
		ctxLogger = ctxLogger.WithString(attrs[i], attrs[i+1])
	}
	return ctxLogger
}
