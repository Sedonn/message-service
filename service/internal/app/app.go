package app

import (
	"log/slog"

	restapp "github.com/sedonn/message-service/internal/app/rest"
	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/repository/postgresql"
	"github.com/sedonn/message-service/internal/services/message"
)

// App это микросервис сообщений.
type App struct {
	RESTApp *restapp.App
}

// New создает новый микросервис сообщений.
func New(log *slog.Logger, cfg *config.Config) *App {
	const op = "app.New"

	repository, err := postgresql.New(cfg)
	if err != nil {
		panic(err)
	}
	log.Info("database connected", slog.String("op", op), slog.String("database", cfg.DB.Database))

	messageService := message.New(log, repository, repository)

	restApp := restapp.New(log, cfg.REST.Port, messageService)

	return &App{
		RESTApp: restApp,
	}
}
