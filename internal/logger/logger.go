package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	log := debugLogger(zerolog.DebugLevel)
	zerolog.DefaultContextLogger = &log
}

func debugLogger(level zerolog.Level) zerolog.Logger {
	return zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}

func prodLogger(level zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stderr).
		Level(level).
		With().
		Timestamp().
		Logger()
}
