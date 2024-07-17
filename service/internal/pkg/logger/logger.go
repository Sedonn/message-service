package logger

import (
	"log/slog"
	"os"

	"github.com/sedonn/message-service/internal/config"
	gormlog "gorm.io/gorm/logger"
)

// New создает и настраивает объект логгера на основе типа текущего окружения.
func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// New создает и настраивает объект логгера GORM на основе типа текущего окружения.
func NewGORMLogger(env string) gormlog.Interface {
	var level gormlog.LogLevel

	switch env {
	case config.EnvLocal:
		level = gormlog.Info
	case config.EnvProduction:
		level = gormlog.Silent
	}

	return gormlog.Default.LogMode(level)
}