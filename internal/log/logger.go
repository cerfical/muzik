package log

import (
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
	l.log(FatalLevel, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *Logger) log(lvl Level, msg string) {
	l.logger.WithLevel(zerolog.Level(lvl)).Msg(msg)
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
