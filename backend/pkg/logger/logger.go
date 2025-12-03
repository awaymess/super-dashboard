package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Init initializes the logger
func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	log = zerolog.New(os.Stdout).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Get returns the logger instance
func Get() zerolog.Logger {
	return log
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Error logs an error message
func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn().Msg(msg)
}
