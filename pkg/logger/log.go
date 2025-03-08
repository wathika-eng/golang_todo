package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(dsn string) {
	// upm(dsn)
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(Logger)
}
