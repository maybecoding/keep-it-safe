// Package logger - package for logging.
package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// lg Общая переменная для логирования будет доступна всему коду
// Не лучшее решение, но самое простое.
var lg *zerolog.Logger

// Init - init of logger variable using passed log level, for values:
// fatal error info debug
// if value not from list above sets debug log level.
func Init(level string, useFile bool) error {
	var zl zerolog.Logger
	if useFile {
		zl = zerolog.New(os.Stderr).With().Timestamp().Logger()
		file, err := os.Create("log" + time.Now().String() + ".log")
		if err != nil {
			return fmt.Errorf("logger - Init - error due create log file %w", err)
		}
		zl = log.Output(file)
	} else {
		zl = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	switch level {
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zl.Debug().Msg("passed wrong error level")
	}
	zl.Debug().Str("log level", level).Msg("log initialized")

	lg = &zl
	return nil
}

// Fatal - starts a new message with fatal level. The os.Exit(1) function is called by the Msg method,
// which terminates the program immediately
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return lg.Fatal()
}

// Error starts a new message with error level.
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return lg.Error()
}

// Info starts a new message with info level.
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return lg.Info()
}

// Debug starts a new message with debug level.
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return lg.Debug()
}
