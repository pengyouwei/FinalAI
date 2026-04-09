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

// func NewLogger(lv string) *slog.Logger {
// 	level := parseLevel(lv)

// 	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
// 		AddSource:  true,
// 		Level:      level,
// 		TimeFormat: time.RFC3339,
// 	})

// 	logger := slog.New(consoleHandler)

// 	return logger
// }

// func parseLevel(level string) slog.Level {
// 	switch strings.ToLower(level) {
// 	case "debug":
// 		return slog.LevelDebug
// 	case "info":
// 		return slog.LevelInfo
// 	case "warn", "warning":
// 		return slog.LevelWarn
// 	case "error":
// 		return slog.LevelError
// 	default:
// 		return slog.LevelInfo
// 	}
// }
