package logger

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New() *Logger {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &Logger{l}
}
