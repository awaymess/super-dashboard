package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the logger with settings from environment variables.
// LOG_LEVEL env var can be: debug, info, warn, error, fatal, panic (default: info)
// ENV=production enables JSON output; otherwise uses pretty console output.
func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level from environment
	level := getLogLevel()
	zerolog.SetGlobalLevel(level)

	// Pretty print in development
	if os.Getenv("ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().Str("level", level.String()).Msg("Logger initialized")
}

// getLogLevel parses the LOG_LEVEL environment variable.
func getLogLevel() zerolog.Level {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
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

// WithField returns a logger with an additional field.
func WithField(key string, value interface{}) zerolog.Logger {
	return log.With().Interface(key, value).Logger()
}
