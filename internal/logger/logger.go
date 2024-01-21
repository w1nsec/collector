package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Initialize(level string) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvl := selectLevel(level)
	zerolog.SetGlobalLevel(lvl)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	})
	return nil
}

func selectLevel(level string) zerolog.Level {
	level = strings.ToTitle(level)
	switch level {
	// INFO
	case "INF":
		return zerolog.InfoLevel
	case "INFO":
		return zerolog.InfoLevel

	// DEBUG
	case "DBG":
		return zerolog.DebugLevel
	case "DEBUG":
		return zerolog.DebugLevel

	// Warning
	case "WRN":
		return zerolog.WarnLevel
	case "WARNING":
		return zerolog.WarnLevel
	case "WARN":
		return zerolog.WarnLevel

	// Error level
	case "ERR":
		return zerolog.ErrorLevel
	case "ERROR":
		return zerolog.ErrorLevel

	default:
		return zerolog.DebugLevel
	}
}
