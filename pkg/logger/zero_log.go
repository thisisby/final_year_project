package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var ZeroLogger zerolog.Logger

func InitZeroLogger() {
	ZeroLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}
