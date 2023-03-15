package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func New(levelStr string) *zerolog.Logger {
	var level zerolog.Level
	switch levelStr {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.DebugLevel
	}
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	log := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf(" %s:%d ", file, line)
	}
	return &log
}
