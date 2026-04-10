package mylogger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String("time",
					a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
